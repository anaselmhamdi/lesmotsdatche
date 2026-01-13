// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'puzzle_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_$CluesImpl _$$CluesImplFromJson(Map<String, dynamic> json) => _$CluesImpl(
      across: (json['across'] as List<dynamic>)
          .map((e) => Clue.fromJson(e as Map<String, dynamic>))
          .toList(),
      down: (json['down'] as List<dynamic>)
          .map((e) => Clue.fromJson(e as Map<String, dynamic>))
          .toList(),
    );

Map<String, dynamic> _$$CluesImplToJson(_$CluesImpl instance) =>
    <String, dynamic>{
      'across': instance.across,
      'down': instance.down,
    };

_$MetadataImpl _$$MetadataImplFromJson(Map<String, dynamic> json) =>
    _$MetadataImpl(
      themeTags: (json['theme_tags'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList(),
      referenceTags: (json['reference_tags'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList(),
      notes: json['notes'] as String?,
      freshnessScore: (json['freshness_score'] as num?)?.toInt(),
    );

Map<String, dynamic> _$$MetadataImplToJson(_$MetadataImpl instance) =>
    <String, dynamic>{
      'theme_tags': instance.themeTags,
      'reference_tags': instance.referenceTags,
      'notes': instance.notes,
      'freshness_score': instance.freshnessScore,
    };

_$PuzzleImpl _$$PuzzleImplFromJson(Map<String, dynamic> json) => _$PuzzleImpl(
      id: json['id'] as String,
      date: json['date'] as String,
      language: json['language'] as String,
      title: json['title'] as String,
      author: json['author'] as String,
      difficulty: (json['difficulty'] as num).toInt(),
      status: json['status'] as String,
      grid: (json['grid'] as List<dynamic>)
          .map((e) => (e as List<dynamic>)
              .map((e) => Cell.fromJson(e as Map<String, dynamic>))
              .toList())
          .toList(),
      clues: Clues.fromJson(json['clues'] as Map<String, dynamic>),
      metadata: json['metadata'] == null
          ? null
          : Metadata.fromJson(json['metadata'] as Map<String, dynamic>),
      createdAt: DateTime.parse(json['created_at'] as String),
      publishedAt: json['published_at'] == null
          ? null
          : DateTime.parse(json['published_at'] as String),
    );

Map<String, dynamic> _$$PuzzleImplToJson(_$PuzzleImpl instance) =>
    <String, dynamic>{
      'id': instance.id,
      'date': instance.date,
      'language': instance.language,
      'title': instance.title,
      'author': instance.author,
      'difficulty': instance.difficulty,
      'status': instance.status,
      'grid': instance.grid,
      'clues': instance.clues,
      'metadata': instance.metadata,
      'created_at': instance.createdAt.toIso8601String(),
      'published_at': instance.publishedAt?.toIso8601String(),
    };

_$PuzzleSummaryImpl _$$PuzzleSummaryImplFromJson(Map<String, dynamic> json) =>
    _$PuzzleSummaryImpl(
      id: json['id'] as String,
      date: json['date'] as String,
      language: json['language'] as String,
      title: json['title'] as String,
      author: json['author'] as String,
      difficulty: (json['difficulty'] as num).toInt(),
      status: json['status'] as String,
    );

Map<String, dynamic> _$$PuzzleSummaryImplToJson(_$PuzzleSummaryImpl instance) =>
    <String, dynamic>{
      'id': instance.id,
      'date': instance.date,
      'language': instance.language,
      'title': instance.title,
      'author': instance.author,
      'difficulty': instance.difficulty,
      'status': instance.status,
    };
