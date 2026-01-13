// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'cell_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_$CellImpl _$$CellImplFromJson(Map<String, dynamic> json) => _$CellImpl(
      type: json['type'] as String,
      solution: json['solution'] as String?,
      number: (json['number'] as num?)?.toInt(),
      clueAcross: json['clue_across'] as String?,
      clueDown: json['clue_down'] as String?,
    );

Map<String, dynamic> _$$CellImplToJson(_$CellImpl instance) =>
    <String, dynamic>{
      'type': instance.type,
      'solution': instance.solution,
      'number': instance.number,
      'clue_across': instance.clueAcross,
      'clue_down': instance.clueDown,
    };
