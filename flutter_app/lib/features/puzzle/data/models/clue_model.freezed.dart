// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'clue_model.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
    'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more information: https://github.com/rrousselGit/freezed#adding-getters-and-methods-to-our-models');

Clue _$ClueFromJson(Map<String, dynamic> json) {
  return _Clue.fromJson(json);
}

/// @nodoc
mixin _$Clue {
  /// Unique identifier.
  String get id => throw _privateConstructorUsedError;

  /// Direction: "across" or "down".
  String get direction => throw _privateConstructorUsedError;

  /// Clue number.
  int get number => throw _privateConstructorUsedError;

  /// The clue prompt/hint.
  String get prompt => throw _privateConstructorUsedError;

  /// Normalized answer (A-Z, no spaces).
  String get answer => throw _privateConstructorUsedError;

  /// Original answer with spaces/hyphens.
  @JsonKey(name: 'original_answer')
  String? get originalAnswer => throw _privateConstructorUsedError;

  /// Starting position in the grid.
  Position get start => throw _privateConstructorUsedError;

  /// Length of the answer.
  int get length => throw _privateConstructorUsedError;

  /// Reference tags for the clue.
  @JsonKey(name: 'reference_tags')
  List<String>? get referenceTags => throw _privateConstructorUsedError;

  /// Year range for references.
  @JsonKey(name: 'reference_year_range')
  List<int>? get referenceYearRange => throw _privateConstructorUsedError;

  /// Difficulty level (1-5).
  int? get difficulty => throw _privateConstructorUsedError;

  /// Notes about ambiguity.
  @JsonKey(name: 'ambiguity_notes')
  String? get ambiguityNotes => throw _privateConstructorUsedError;

  /// Serializes this Clue to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of Clue
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $ClueCopyWith<Clue> get copyWith => throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $ClueCopyWith<$Res> {
  factory $ClueCopyWith(Clue value, $Res Function(Clue) then) =
      _$ClueCopyWithImpl<$Res, Clue>;
  @useResult
  $Res call(
      {String id,
      String direction,
      int number,
      String prompt,
      String answer,
      @JsonKey(name: 'original_answer') String? originalAnswer,
      Position start,
      int length,
      @JsonKey(name: 'reference_tags') List<String>? referenceTags,
      @JsonKey(name: 'reference_year_range') List<int>? referenceYearRange,
      int? difficulty,
      @JsonKey(name: 'ambiguity_notes') String? ambiguityNotes});

  $PositionCopyWith<$Res> get start;
}

/// @nodoc
class _$ClueCopyWithImpl<$Res, $Val extends Clue>
    implements $ClueCopyWith<$Res> {
  _$ClueCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of Clue
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? direction = null,
    Object? number = null,
    Object? prompt = null,
    Object? answer = null,
    Object? originalAnswer = freezed,
    Object? start = null,
    Object? length = null,
    Object? referenceTags = freezed,
    Object? referenceYearRange = freezed,
    Object? difficulty = freezed,
    Object? ambiguityNotes = freezed,
  }) {
    return _then(_value.copyWith(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as String,
      direction: null == direction
          ? _value.direction
          : direction // ignore: cast_nullable_to_non_nullable
              as String,
      number: null == number
          ? _value.number
          : number // ignore: cast_nullable_to_non_nullable
              as int,
      prompt: null == prompt
          ? _value.prompt
          : prompt // ignore: cast_nullable_to_non_nullable
              as String,
      answer: null == answer
          ? _value.answer
          : answer // ignore: cast_nullable_to_non_nullable
              as String,
      originalAnswer: freezed == originalAnswer
          ? _value.originalAnswer
          : originalAnswer // ignore: cast_nullable_to_non_nullable
              as String?,
      start: null == start
          ? _value.start
          : start // ignore: cast_nullable_to_non_nullable
              as Position,
      length: null == length
          ? _value.length
          : length // ignore: cast_nullable_to_non_nullable
              as int,
      referenceTags: freezed == referenceTags
          ? _value.referenceTags
          : referenceTags // ignore: cast_nullable_to_non_nullable
              as List<String>?,
      referenceYearRange: freezed == referenceYearRange
          ? _value.referenceYearRange
          : referenceYearRange // ignore: cast_nullable_to_non_nullable
              as List<int>?,
      difficulty: freezed == difficulty
          ? _value.difficulty
          : difficulty // ignore: cast_nullable_to_non_nullable
              as int?,
      ambiguityNotes: freezed == ambiguityNotes
          ? _value.ambiguityNotes
          : ambiguityNotes // ignore: cast_nullable_to_non_nullable
              as String?,
    ) as $Val);
  }

  /// Create a copy of Clue
  /// with the given fields replaced by the non-null parameter values.
  @override
  @pragma('vm:prefer-inline')
  $PositionCopyWith<$Res> get start {
    return $PositionCopyWith<$Res>(_value.start, (value) {
      return _then(_value.copyWith(start: value) as $Val);
    });
  }
}

/// @nodoc
abstract class _$$ClueImplCopyWith<$Res> implements $ClueCopyWith<$Res> {
  factory _$$ClueImplCopyWith(
          _$ClueImpl value, $Res Function(_$ClueImpl) then) =
      __$$ClueImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {String id,
      String direction,
      int number,
      String prompt,
      String answer,
      @JsonKey(name: 'original_answer') String? originalAnswer,
      Position start,
      int length,
      @JsonKey(name: 'reference_tags') List<String>? referenceTags,
      @JsonKey(name: 'reference_year_range') List<int>? referenceYearRange,
      int? difficulty,
      @JsonKey(name: 'ambiguity_notes') String? ambiguityNotes});

  @override
  $PositionCopyWith<$Res> get start;
}

/// @nodoc
class __$$ClueImplCopyWithImpl<$Res>
    extends _$ClueCopyWithImpl<$Res, _$ClueImpl>
    implements _$$ClueImplCopyWith<$Res> {
  __$$ClueImplCopyWithImpl(_$ClueImpl _value, $Res Function(_$ClueImpl) _then)
      : super(_value, _then);

  /// Create a copy of Clue
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? direction = null,
    Object? number = null,
    Object? prompt = null,
    Object? answer = null,
    Object? originalAnswer = freezed,
    Object? start = null,
    Object? length = null,
    Object? referenceTags = freezed,
    Object? referenceYearRange = freezed,
    Object? difficulty = freezed,
    Object? ambiguityNotes = freezed,
  }) {
    return _then(_$ClueImpl(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as String,
      direction: null == direction
          ? _value.direction
          : direction // ignore: cast_nullable_to_non_nullable
              as String,
      number: null == number
          ? _value.number
          : number // ignore: cast_nullable_to_non_nullable
              as int,
      prompt: null == prompt
          ? _value.prompt
          : prompt // ignore: cast_nullable_to_non_nullable
              as String,
      answer: null == answer
          ? _value.answer
          : answer // ignore: cast_nullable_to_non_nullable
              as String,
      originalAnswer: freezed == originalAnswer
          ? _value.originalAnswer
          : originalAnswer // ignore: cast_nullable_to_non_nullable
              as String?,
      start: null == start
          ? _value.start
          : start // ignore: cast_nullable_to_non_nullable
              as Position,
      length: null == length
          ? _value.length
          : length // ignore: cast_nullable_to_non_nullable
              as int,
      referenceTags: freezed == referenceTags
          ? _value._referenceTags
          : referenceTags // ignore: cast_nullable_to_non_nullable
              as List<String>?,
      referenceYearRange: freezed == referenceYearRange
          ? _value._referenceYearRange
          : referenceYearRange // ignore: cast_nullable_to_non_nullable
              as List<int>?,
      difficulty: freezed == difficulty
          ? _value.difficulty
          : difficulty // ignore: cast_nullable_to_non_nullable
              as int?,
      ambiguityNotes: freezed == ambiguityNotes
          ? _value.ambiguityNotes
          : ambiguityNotes // ignore: cast_nullable_to_non_nullable
              as String?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$ClueImpl extends _Clue {
  const _$ClueImpl(
      {required this.id,
      required this.direction,
      required this.number,
      required this.prompt,
      required this.answer,
      @JsonKey(name: 'original_answer') this.originalAnswer,
      required this.start,
      required this.length,
      @JsonKey(name: 'reference_tags') final List<String>? referenceTags,
      @JsonKey(name: 'reference_year_range')
      final List<int>? referenceYearRange,
      this.difficulty,
      @JsonKey(name: 'ambiguity_notes') this.ambiguityNotes})
      : _referenceTags = referenceTags,
        _referenceYearRange = referenceYearRange,
        super._();

  factory _$ClueImpl.fromJson(Map<String, dynamic> json) =>
      _$$ClueImplFromJson(json);

  /// Unique identifier.
  @override
  final String id;

  /// Direction: "across" or "down".
  @override
  final String direction;

  /// Clue number.
  @override
  final int number;

  /// The clue prompt/hint.
  @override
  final String prompt;

  /// Normalized answer (A-Z, no spaces).
  @override
  final String answer;

  /// Original answer with spaces/hyphens.
  @override
  @JsonKey(name: 'original_answer')
  final String? originalAnswer;

  /// Starting position in the grid.
  @override
  final Position start;

  /// Length of the answer.
  @override
  final int length;

  /// Reference tags for the clue.
  final List<String>? _referenceTags;

  /// Reference tags for the clue.
  @override
  @JsonKey(name: 'reference_tags')
  List<String>? get referenceTags {
    final value = _referenceTags;
    if (value == null) return null;
    if (_referenceTags is EqualUnmodifiableListView) return _referenceTags;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableListView(value);
  }

  /// Year range for references.
  final List<int>? _referenceYearRange;

  /// Year range for references.
  @override
  @JsonKey(name: 'reference_year_range')
  List<int>? get referenceYearRange {
    final value = _referenceYearRange;
    if (value == null) return null;
    if (_referenceYearRange is EqualUnmodifiableListView)
      return _referenceYearRange;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableListView(value);
  }

  /// Difficulty level (1-5).
  @override
  final int? difficulty;

  /// Notes about ambiguity.
  @override
  @JsonKey(name: 'ambiguity_notes')
  final String? ambiguityNotes;

  @override
  String toString() {
    return 'Clue(id: $id, direction: $direction, number: $number, prompt: $prompt, answer: $answer, originalAnswer: $originalAnswer, start: $start, length: $length, referenceTags: $referenceTags, referenceYearRange: $referenceYearRange, difficulty: $difficulty, ambiguityNotes: $ambiguityNotes)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$ClueImpl &&
            (identical(other.id, id) || other.id == id) &&
            (identical(other.direction, direction) ||
                other.direction == direction) &&
            (identical(other.number, number) || other.number == number) &&
            (identical(other.prompt, prompt) || other.prompt == prompt) &&
            (identical(other.answer, answer) || other.answer == answer) &&
            (identical(other.originalAnswer, originalAnswer) ||
                other.originalAnswer == originalAnswer) &&
            (identical(other.start, start) || other.start == start) &&
            (identical(other.length, length) || other.length == length) &&
            const DeepCollectionEquality()
                .equals(other._referenceTags, _referenceTags) &&
            const DeepCollectionEquality()
                .equals(other._referenceYearRange, _referenceYearRange) &&
            (identical(other.difficulty, difficulty) ||
                other.difficulty == difficulty) &&
            (identical(other.ambiguityNotes, ambiguityNotes) ||
                other.ambiguityNotes == ambiguityNotes));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(
      runtimeType,
      id,
      direction,
      number,
      prompt,
      answer,
      originalAnswer,
      start,
      length,
      const DeepCollectionEquality().hash(_referenceTags),
      const DeepCollectionEquality().hash(_referenceYearRange),
      difficulty,
      ambiguityNotes);

  /// Create a copy of Clue
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$ClueImplCopyWith<_$ClueImpl> get copyWith =>
      __$$ClueImplCopyWithImpl<_$ClueImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$ClueImplToJson(
      this,
    );
  }
}

abstract class _Clue extends Clue {
  const factory _Clue(
          {required final String id,
          required final String direction,
          required final int number,
          required final String prompt,
          required final String answer,
          @JsonKey(name: 'original_answer') final String? originalAnswer,
          required final Position start,
          required final int length,
          @JsonKey(name: 'reference_tags') final List<String>? referenceTags,
          @JsonKey(name: 'reference_year_range')
          final List<int>? referenceYearRange,
          final int? difficulty,
          @JsonKey(name: 'ambiguity_notes') final String? ambiguityNotes}) =
      _$ClueImpl;
  const _Clue._() : super._();

  factory _Clue.fromJson(Map<String, dynamic> json) = _$ClueImpl.fromJson;

  /// Unique identifier.
  @override
  String get id;

  /// Direction: "across" or "down".
  @override
  String get direction;

  /// Clue number.
  @override
  int get number;

  /// The clue prompt/hint.
  @override
  String get prompt;

  /// Normalized answer (A-Z, no spaces).
  @override
  String get answer;

  /// Original answer with spaces/hyphens.
  @override
  @JsonKey(name: 'original_answer')
  String? get originalAnswer;

  /// Starting position in the grid.
  @override
  Position get start;

  /// Length of the answer.
  @override
  int get length;

  /// Reference tags for the clue.
  @override
  @JsonKey(name: 'reference_tags')
  List<String>? get referenceTags;

  /// Year range for references.
  @override
  @JsonKey(name: 'reference_year_range')
  List<int>? get referenceYearRange;

  /// Difficulty level (1-5).
  @override
  int? get difficulty;

  /// Notes about ambiguity.
  @override
  @JsonKey(name: 'ambiguity_notes')
  String? get ambiguityNotes;

  /// Create a copy of Clue
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$ClueImplCopyWith<_$ClueImpl> get copyWith =>
      throw _privateConstructorUsedError;
}
