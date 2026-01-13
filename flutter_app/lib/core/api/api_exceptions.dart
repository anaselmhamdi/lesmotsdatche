/// Base class for API exceptions.
sealed class ApiException implements Exception {
  final String message;
  final int? statusCode;

  const ApiException(this.message, {this.statusCode});

  @override
  String toString() => 'ApiException: $message (status: $statusCode)';
}

/// Network error (no connection, timeout, etc.).
class NetworkException extends ApiException {
  const NetworkException(super.message);
}

/// Server returned an error response.
class ServerException extends ApiException {
  const ServerException(super.message, {super.statusCode});
}

/// Resource not found (404).
class NotFoundException extends ApiException {
  const NotFoundException([String message = 'Resource not found'])
      : super(message, statusCode: 404);
}

/// Unauthorized (401).
class UnauthorizedException extends ApiException {
  const UnauthorizedException([String message = 'Unauthorized'])
      : super(message, statusCode: 401);
}
