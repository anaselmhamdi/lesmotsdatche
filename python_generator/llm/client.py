"""OpenAI LLM client with validation and retry logic."""

import json
import re
from dataclasses import dataclass, field
from typing import Any, TypeVar

from openai import OpenAI
from pydantic import BaseModel, ValidationError


class LLMError(Exception):
    """LLM-related error."""
    pass


@dataclass
class LLMConfig:
    """LLM client configuration."""

    max_retries: int = 3
    default_temperature: float = 0.7
    default_max_tokens: int = 2048
    model: str = "gpt-4o"
    repair_prompt: str = """The previous response was invalid JSON or didn't match the required schema.
Error: {error}
Previous response: {response}

Please provide a corrected response that is valid JSON matching the schema."""


@dataclass
class Trace:
    """Records an LLM interaction for debugging."""

    prompt: str
    response: str
    error: str = ""
    attempt: int = 1


T = TypeVar("T", bound=BaseModel)


class LLMClient:
    """OpenAI client with JSON validation and retry logic."""

    def __init__(self, api_key: str | None = None, config: LLMConfig | None = None):
        self.client = OpenAI(api_key=api_key)
        self.config = config or LLMConfig()
        self.traces: list[Trace] = []

    def complete(
        self,
        prompt: str,
        response_model: type[T],
        system_prompt: str | None = None,
        temperature: float | None = None,
        max_tokens: int | None = None,
    ) -> T:
        """Send a completion request and parse the response into a Pydantic model.

        Args:
            prompt: User prompt
            response_model: Pydantic model class for response parsing
            system_prompt: Optional system prompt
            temperature: Sampling temperature (default from config)
            max_tokens: Maximum tokens (default from config)

        Returns:
            Parsed response as the specified Pydantic model

        Raises:
            LLMError: If max retries exceeded or other LLM error
        """
        temp = temperature if temperature is not None else self.config.default_temperature
        tokens = max_tokens if max_tokens is not None else self.config.default_max_tokens

        messages = []
        if system_prompt:
            messages.append({"role": "system", "content": system_prompt})
        messages.append({"role": "user", "content": prompt})

        last_error = None
        current_prompt = prompt

        for attempt in range(1, self.config.max_retries + 1):
            try:
                # Update user message with repair prompt if retrying
                if attempt > 1:
                    messages[-1] = {"role": "user", "content": current_prompt}

                response = self.client.chat.completions.create(
                    model=self.config.model,
                    messages=messages,
                    temperature=temp,
                    max_tokens=tokens,
                )

                content = response.choices[0].message.content or ""
                self.traces.append(Trace(
                    prompt=current_prompt,
                    response=content,
                    attempt=attempt,
                ))

                if not content.strip():
                    last_error = "Empty response from LLM"
                    current_prompt = self.config.repair_prompt.format(
                        error=last_error,
                        response="(empty)",
                    )
                    continue

                # Extract JSON from response (handle markdown code blocks)
                json_content = _extract_json(content)

                # Parse with Pydantic
                try:
                    return response_model.model_validate_json(json_content)
                except ValidationError as e:
                    last_error = f"Validation error: {e}"
                    current_prompt = self.config.repair_prompt.format(
                        error=last_error,
                        response=_truncate(content, 500),
                    )
                    self.traces[-1].error = last_error
                    continue

            except json.JSONDecodeError as e:
                last_error = f"JSON parse error: {e}"
                current_prompt = self.config.repair_prompt.format(
                    error=last_error,
                    response=_truncate(content, 500),
                )
                self.traces.append(Trace(
                    prompt=current_prompt,
                    response=content if 'content' in dir() else "",
                    error=last_error,
                    attempt=attempt,
                ))
                continue

        raise LLMError(f"Max retries exceeded: {last_error}")

    def complete_raw(
        self,
        prompt: str,
        system_prompt: str | None = None,
        temperature: float | None = None,
        max_tokens: int | None = None,
    ) -> str:
        """Send a completion request and return raw text response.

        Args:
            prompt: User prompt
            system_prompt: Optional system prompt
            temperature: Sampling temperature (default from config)
            max_tokens: Maximum tokens (default from config)

        Returns:
            Raw text response

        Raises:
            LLMError: If LLM request fails
        """
        temp = temperature if temperature is not None else self.config.default_temperature
        tokens = max_tokens if max_tokens is not None else self.config.default_max_tokens

        messages = []
        if system_prompt:
            messages.append({"role": "system", "content": system_prompt})
        messages.append({"role": "user", "content": prompt})

        try:
            response = self.client.chat.completions.create(
                model=self.config.model,
                messages=messages,
                temperature=temp,
                max_tokens=tokens,
            )
            content = response.choices[0].message.content or ""
            self.traces.append(Trace(prompt=prompt, response=content))
            return content
        except Exception as e:
            raise LLMError(f"LLM request failed: {e}")

    def clear_traces(self):
        """Clear recorded traces."""
        self.traces = []


def _extract_json(content: str) -> str:
    """Extract JSON from a response that might be wrapped in markdown."""
    content = content.strip()

    # Check for markdown code blocks
    if content.startswith("```"):
        # Find the end of the first line (language specifier)
        first_newline = content.find("\n")
        if first_newline > 0:
            content = content[first_newline + 1:]
        # Remove trailing ```
        last_fence = content.rfind("```")
        if last_fence > 0:
            content = content[:last_fence]

    return content.strip()


def _truncate(s: str, max_len: int) -> str:
    """Truncate string to max length."""
    if len(s) <= max_len:
        return s
    return s[:max_len] + "..."
