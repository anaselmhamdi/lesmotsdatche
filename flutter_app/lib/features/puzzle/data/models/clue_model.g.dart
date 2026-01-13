// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'clue_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_$ClueImpl _$$ClueImplFromJson(Map<String, dynamic> json) => _$ClueImpl(
      id: json['id'] as String,
      direction: json['direction'] as String,
      number: (json['number'] as num).toInt(),
      prompt: json['prompt'] as String,
      answer: json['answer'] as String,
      originalAnswer: json['original_answer'] as String?,
      start: Position.fromJson(json['start'] as Map<String, dynamic>),
      length: (json['length'] as num).toInt(),
      referenceTags: (json['reference_tags'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList(),
      referenceYearRange: (json['reference_year_range'] as List<dynamic>?)
          ?.map((e) => (e as num).toInt())
          .toList(),
      difficulty: (json['difficulty'] as num?)?.toInt(),
      ambiguityNotes: json['ambiguity_notes'] as String?,
    );

Map<String, dynamic> _$$ClueImplToJson(_$ClueImpl instance) =>
    <String, dynamic>{
      'id': instance.id,
      'direction': instance.direction,
      'number': instance.number,
      'prompt': instance.prompt,
      'answer': instance.answer,
      'original_answer': instance.originalAnswer,
      'start': instance.start,
      'length': instance.length,
      'reference_tags': instance.referenceTags,
      'reference_year_range': instance.referenceYearRange,
      'difficulty': instance.difficulty,
      'ambiguity_notes': instance.ambiguityNotes,
    };
