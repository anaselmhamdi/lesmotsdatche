// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'puzzle_model.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
    'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more information: https://github.com/rrousselGit/freezed#adding-getters-and-methods-to-our-models');

Clues _$CluesFromJson(Map<String, dynamic> json) {
  return _Clues.fromJson(json);
}

/// @nodoc
mixin _$Clues {
  List<Clue> get across => throw _privateConstructorUsedError;
  List<Clue> get down => throw _privateConstructorUsedError;

  /// Serializes this Clues to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of Clues
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $CluesCopyWith<Clues> get copyWith => throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $CluesCopyWith<$Res> {
  factory $CluesCopyWith(Clues value, $Res Function(Clues) then) =
      _$CluesCopyWithImpl<$Res, Clues>;
  @useResult
  $Res call({List<Clue> across, List<Clue> down});
}

/// @nodoc
class _$CluesCopyWithImpl<$Res, $Val extends Clues>
    implements $CluesCopyWith<$Res> {
  _$CluesCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of Clues
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? across = null,
    Object? down = null,
  }) {
    return _then(_value.copyWith(
      across: null == across
          ? _value.across
          : across // ignore: cast_nullable_to_non_nullable
              as List<Clue>,
      down: null == down
          ? _value.down
          : down // ignore: cast_nullable_to_non_nullable
              as List<Clue>,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$CluesImplCopyWith<$Res> implements $CluesCopyWith<$Res> {
  factory _$$CluesImplCopyWith(
          _$CluesImpl value, $Res Function(_$CluesImpl) then) =
      __$$CluesImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({List<Clue> across, List<Clue> down});
}

/// @nodoc
class __$$CluesImplCopyWithImpl<$Res>
    extends _$CluesCopyWithImpl<$Res, _$CluesImpl>
    implements _$$CluesImplCopyWith<$Res> {
  __$$CluesImplCopyWithImpl(
      _$CluesImpl _value, $Res Function(_$CluesImpl) _then)
      : super(_value, _then);

  /// Create a copy of Clues
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? across = null,
    Object? down = null,
  }) {
    return _then(_$CluesImpl(
      across: null == across
          ? _value._across
          : across // ignore: cast_nullable_to_non_nullable
              as List<Clue>,
      down: null == down
          ? _value._down
          : down // ignore: cast_nullable_to_non_nullable
              as List<Clue>,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$CluesImpl implements _Clues {
  const _$CluesImpl(
      {required final List<Clue> across, required final List<Clue> down})
      : _across = across,
        _down = down;

  factory _$CluesImpl.fromJson(Map<String, dynamic> json) =>
      _$$CluesImplFromJson(json);

  final List<Clue> _across;
  @override
  List<Clue> get across {
    if (_across is EqualUnmodifiableListView) return _across;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableListView(_across);
  }

  final List<Clue> _down;
  @override
  List<Clue> get down {
    if (_down is EqualUnmodifiableListView) return _down;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableListView(_down);
  }

  @override
  String toString() {
    return 'Clues(across: $across, down: $down)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$CluesImpl &&
            const DeepCollectionEquality().equals(other._across, _across) &&
            const DeepCollectionEquality().equals(other._down, _down));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(
      runtimeType,
      const DeepCollectionEquality().hash(_across),
      const DeepCollectionEquality().hash(_down));

  /// Create a copy of Clues
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$CluesImplCopyWith<_$CluesImpl> get copyWith =>
      __$$CluesImplCopyWithImpl<_$CluesImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$CluesImplToJson(
      this,
    );
  }
}

abstract class _Clues implements Clues {
  const factory _Clues(
      {required final List<Clue> across,
      required final List<Clue> down}) = _$CluesImpl;

  factory _Clues.fromJson(Map<String, dynamic> json) = _$CluesImpl.fromJson;

  @override
  List<Clue> get across;
  @override
  List<Clue> get down;

  /// Create a copy of Clues
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$CluesImplCopyWith<_$CluesImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

Metadata _$MetadataFromJson(Map<String, dynamic> json) {
  return _Metadata.fromJson(json);
}

/// @nodoc
mixin _$Metadata {
  @JsonKey(name: 'theme_tags')
  List<String>? get themeTags => throw _privateConstructorUsedError;
  @JsonKey(name: 'reference_tags')
  List<String>? get referenceTags => throw _privateConstructorUsedError;
  String? get notes => throw _privateConstructorUsedError;
  @JsonKey(name: 'freshness_score')
  int? get freshnessScore => throw _privateConstructorUsedError;

  /// Serializes this Metadata to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of Metadata
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $MetadataCopyWith<Metadata> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $MetadataCopyWith<$Res> {
  factory $MetadataCopyWith(Metadata value, $Res Function(Metadata) then) =
      _$MetadataCopyWithImpl<$Res, Metadata>;
  @useResult
  $Res call(
      {@JsonKey(name: 'theme_tags') List<String>? themeTags,
      @JsonKey(name: 'reference_tags') List<String>? referenceTags,
      String? notes,
      @JsonKey(name: 'freshness_score') int? freshnessScore});
}

/// @nodoc
class _$MetadataCopyWithImpl<$Res, $Val extends Metadata>
    implements $MetadataCopyWith<$Res> {
  _$MetadataCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of Metadata
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? themeTags = freezed,
    Object? referenceTags = freezed,
    Object? notes = freezed,
    Object? freshnessScore = freezed,
  }) {
    return _then(_value.copyWith(
      themeTags: freezed == themeTags
          ? _value.themeTags
          : themeTags // ignore: cast_nullable_to_non_nullable
              as List<String>?,
      referenceTags: freezed == referenceTags
          ? _value.referenceTags
          : referenceTags // ignore: cast_nullable_to_non_nullable
              as List<String>?,
      notes: freezed == notes
          ? _value.notes
          : notes // ignore: cast_nullable_to_non_nullable
              as String?,
      freshnessScore: freezed == freshnessScore
          ? _value.freshnessScore
          : freshnessScore // ignore: cast_nullable_to_non_nullable
              as int?,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$MetadataImplCopyWith<$Res>
    implements $MetadataCopyWith<$Res> {
  factory _$$MetadataImplCopyWith(
          _$MetadataImpl value, $Res Function(_$MetadataImpl) then) =
      __$$MetadataImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {@JsonKey(name: 'theme_tags') List<String>? themeTags,
      @JsonKey(name: 'reference_tags') List<String>? referenceTags,
      String? notes,
      @JsonKey(name: 'freshness_score') int? freshnessScore});
}

/// @nodoc
class __$$MetadataImplCopyWithImpl<$Res>
    extends _$MetadataCopyWithImpl<$Res, _$MetadataImpl>
    implements _$$MetadataImplCopyWith<$Res> {
  __$$MetadataImplCopyWithImpl(
      _$MetadataImpl _value, $Res Function(_$MetadataImpl) _then)
      : super(_value, _then);

  /// Create a copy of Metadata
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? themeTags = freezed,
    Object? referenceTags = freezed,
    Object? notes = freezed,
    Object? freshnessScore = freezed,
  }) {
    return _then(_$MetadataImpl(
      themeTags: freezed == themeTags
          ? _value._themeTags
          : themeTags // ignore: cast_nullable_to_non_nullable
              as List<String>?,
      referenceTags: freezed == referenceTags
          ? _value._referenceTags
          : referenceTags // ignore: cast_nullable_to_non_nullable
              as List<String>?,
      notes: freezed == notes
          ? _value.notes
          : notes // ignore: cast_nullable_to_non_nullable
              as String?,
      freshnessScore: freezed == freshnessScore
          ? _value.freshnessScore
          : freshnessScore // ignore: cast_nullable_to_non_nullable
              as int?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$MetadataImpl implements _Metadata {
  const _$MetadataImpl(
      {@JsonKey(name: 'theme_tags') final List<String>? themeTags,
      @JsonKey(name: 'reference_tags') final List<String>? referenceTags,
      this.notes,
      @JsonKey(name: 'freshness_score') this.freshnessScore})
      : _themeTags = themeTags,
        _referenceTags = referenceTags;

  factory _$MetadataImpl.fromJson(Map<String, dynamic> json) =>
      _$$MetadataImplFromJson(json);

  final List<String>? _themeTags;
  @override
  @JsonKey(name: 'theme_tags')
  List<String>? get themeTags {
    final value = _themeTags;
    if (value == null) return null;
    if (_themeTags is EqualUnmodifiableListView) return _themeTags;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableListView(value);
  }

  final List<String>? _referenceTags;
  @override
  @JsonKey(name: 'reference_tags')
  List<String>? get referenceTags {
    final value = _referenceTags;
    if (value == null) return null;
    if (_referenceTags is EqualUnmodifiableListView) return _referenceTags;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableListView(value);
  }

  @override
  final String? notes;
  @override
  @JsonKey(name: 'freshness_score')
  final int? freshnessScore;

  @override
  String toString() {
    return 'Metadata(themeTags: $themeTags, referenceTags: $referenceTags, notes: $notes, freshnessScore: $freshnessScore)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$MetadataImpl &&
            const DeepCollectionEquality()
                .equals(other._themeTags, _themeTags) &&
            const DeepCollectionEquality()
                .equals(other._referenceTags, _referenceTags) &&
            (identical(other.notes, notes) || other.notes == notes) &&
            (identical(other.freshnessScore, freshnessScore) ||
                other.freshnessScore == freshnessScore));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(
      runtimeType,
      const DeepCollectionEquality().hash(_themeTags),
      const DeepCollectionEquality().hash(_referenceTags),
      notes,
      freshnessScore);

  /// Create a copy of Metadata
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$MetadataImplCopyWith<_$MetadataImpl> get copyWith =>
      __$$MetadataImplCopyWithImpl<_$MetadataImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$MetadataImplToJson(
      this,
    );
  }
}

abstract class _Metadata implements Metadata {
  const factory _Metadata(
          {@JsonKey(name: 'theme_tags') final List<String>? themeTags,
          @JsonKey(name: 'reference_tags') final List<String>? referenceTags,
          final String? notes,
          @JsonKey(name: 'freshness_score') final int? freshnessScore}) =
      _$MetadataImpl;

  factory _Metadata.fromJson(Map<String, dynamic> json) =
      _$MetadataImpl.fromJson;

  @override
  @JsonKey(name: 'theme_tags')
  List<String>? get themeTags;
  @override
  @JsonKey(name: 'reference_tags')
  List<String>? get referenceTags;
  @override
  String? get notes;
  @override
  @JsonKey(name: 'freshness_score')
  int? get freshnessScore;

  /// Create a copy of Metadata
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$MetadataImplCopyWith<_$MetadataImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

Puzzle _$PuzzleFromJson(Map<String, dynamic> json) {
  return _Puzzle.fromJson(json);
}

/// @nodoc
mixin _$Puzzle {
  /// Unique identifier.
  String get id => throw _privateConstructorUsedError;

  /// Publication date (YYYY-MM-DD).
  String get date => throw _privateConstructorUsedError;

  /// Language code ("fr" or "en").
  String get language => throw _privateConstructorUsedError;

  /// Puzzle title.
  String get title => throw _privateConstructorUsedError;

  /// Author name.
  String get author => throw _privateConstructorUsedError;

  /// Difficulty level (1-5).
  int get difficulty => throw _privateConstructorUsedError;

  /// Status ("draft", "published", "archived").
  String get status => throw _privateConstructorUsedError;

  /// 2D grid of cells.
  List<List<Cell>> get grid => throw _privateConstructorUsedError;

  /// Across and down clues.
  Clues get clues => throw _privateConstructorUsedError;

  /// Optional metadata.
  Metadata? get metadata => throw _privateConstructorUsedError;

  /// Creation timestamp.
  @JsonKey(name: 'created_at')
  DateTime get createdAt => throw _privateConstructorUsedError;

  /// Publication timestamp.
  @JsonKey(name: 'published_at')
  DateTime? get publishedAt => throw _privateConstructorUsedError;

  /// Serializes this Puzzle to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of Puzzle
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $PuzzleCopyWith<Puzzle> get copyWith => throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $PuzzleCopyWith<$Res> {
  factory $PuzzleCopyWith(Puzzle value, $Res Function(Puzzle) then) =
      _$PuzzleCopyWithImpl<$Res, Puzzle>;
  @useResult
  $Res call(
      {String id,
      String date,
      String language,
      String title,
      String author,
      int difficulty,
      String status,
      List<List<Cell>> grid,
      Clues clues,
      Metadata? metadata,
      @JsonKey(name: 'created_at') DateTime createdAt,
      @JsonKey(name: 'published_at') DateTime? publishedAt});

  $CluesCopyWith<$Res> get clues;
  $MetadataCopyWith<$Res>? get metadata;
}

/// @nodoc
class _$PuzzleCopyWithImpl<$Res, $Val extends Puzzle>
    implements $PuzzleCopyWith<$Res> {
  _$PuzzleCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of Puzzle
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? date = null,
    Object? language = null,
    Object? title = null,
    Object? author = null,
    Object? difficulty = null,
    Object? status = null,
    Object? grid = null,
    Object? clues = null,
    Object? metadata = freezed,
    Object? createdAt = null,
    Object? publishedAt = freezed,
  }) {
    return _then(_value.copyWith(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as String,
      date: null == date
          ? _value.date
          : date // ignore: cast_nullable_to_non_nullable
              as String,
      language: null == language
          ? _value.language
          : language // ignore: cast_nullable_to_non_nullable
              as String,
      title: null == title
          ? _value.title
          : title // ignore: cast_nullable_to_non_nullable
              as String,
      author: null == author
          ? _value.author
          : author // ignore: cast_nullable_to_non_nullable
              as String,
      difficulty: null == difficulty
          ? _value.difficulty
          : difficulty // ignore: cast_nullable_to_non_nullable
              as int,
      status: null == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as String,
      grid: null == grid
          ? _value.grid
          : grid // ignore: cast_nullable_to_non_nullable
              as List<List<Cell>>,
      clues: null == clues
          ? _value.clues
          : clues // ignore: cast_nullable_to_non_nullable
              as Clues,
      metadata: freezed == metadata
          ? _value.metadata
          : metadata // ignore: cast_nullable_to_non_nullable
              as Metadata?,
      createdAt: null == createdAt
          ? _value.createdAt
          : createdAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
      publishedAt: freezed == publishedAt
          ? _value.publishedAt
          : publishedAt // ignore: cast_nullable_to_non_nullable
              as DateTime?,
    ) as $Val);
  }

  /// Create a copy of Puzzle
  /// with the given fields replaced by the non-null parameter values.
  @override
  @pragma('vm:prefer-inline')
  $CluesCopyWith<$Res> get clues {
    return $CluesCopyWith<$Res>(_value.clues, (value) {
      return _then(_value.copyWith(clues: value) as $Val);
    });
  }

  /// Create a copy of Puzzle
  /// with the given fields replaced by the non-null parameter values.
  @override
  @pragma('vm:prefer-inline')
  $MetadataCopyWith<$Res>? get metadata {
    if (_value.metadata == null) {
      return null;
    }

    return $MetadataCopyWith<$Res>(_value.metadata!, (value) {
      return _then(_value.copyWith(metadata: value) as $Val);
    });
  }
}

/// @nodoc
abstract class _$$PuzzleImplCopyWith<$Res> implements $PuzzleCopyWith<$Res> {
  factory _$$PuzzleImplCopyWith(
          _$PuzzleImpl value, $Res Function(_$PuzzleImpl) then) =
      __$$PuzzleImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {String id,
      String date,
      String language,
      String title,
      String author,
      int difficulty,
      String status,
      List<List<Cell>> grid,
      Clues clues,
      Metadata? metadata,
      @JsonKey(name: 'created_at') DateTime createdAt,
      @JsonKey(name: 'published_at') DateTime? publishedAt});

  @override
  $CluesCopyWith<$Res> get clues;
  @override
  $MetadataCopyWith<$Res>? get metadata;
}

/// @nodoc
class __$$PuzzleImplCopyWithImpl<$Res>
    extends _$PuzzleCopyWithImpl<$Res, _$PuzzleImpl>
    implements _$$PuzzleImplCopyWith<$Res> {
  __$$PuzzleImplCopyWithImpl(
      _$PuzzleImpl _value, $Res Function(_$PuzzleImpl) _then)
      : super(_value, _then);

  /// Create a copy of Puzzle
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? date = null,
    Object? language = null,
    Object? title = null,
    Object? author = null,
    Object? difficulty = null,
    Object? status = null,
    Object? grid = null,
    Object? clues = null,
    Object? metadata = freezed,
    Object? createdAt = null,
    Object? publishedAt = freezed,
  }) {
    return _then(_$PuzzleImpl(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as String,
      date: null == date
          ? _value.date
          : date // ignore: cast_nullable_to_non_nullable
              as String,
      language: null == language
          ? _value.language
          : language // ignore: cast_nullable_to_non_nullable
              as String,
      title: null == title
          ? _value.title
          : title // ignore: cast_nullable_to_non_nullable
              as String,
      author: null == author
          ? _value.author
          : author // ignore: cast_nullable_to_non_nullable
              as String,
      difficulty: null == difficulty
          ? _value.difficulty
          : difficulty // ignore: cast_nullable_to_non_nullable
              as int,
      status: null == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as String,
      grid: null == grid
          ? _value._grid
          : grid // ignore: cast_nullable_to_non_nullable
              as List<List<Cell>>,
      clues: null == clues
          ? _value.clues
          : clues // ignore: cast_nullable_to_non_nullable
              as Clues,
      metadata: freezed == metadata
          ? _value.metadata
          : metadata // ignore: cast_nullable_to_non_nullable
              as Metadata?,
      createdAt: null == createdAt
          ? _value.createdAt
          : createdAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
      publishedAt: freezed == publishedAt
          ? _value.publishedAt
          : publishedAt // ignore: cast_nullable_to_non_nullable
              as DateTime?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$PuzzleImpl extends _Puzzle {
  const _$PuzzleImpl(
      {required this.id,
      required this.date,
      required this.language,
      required this.title,
      required this.author,
      required this.difficulty,
      required this.status,
      required final List<List<Cell>> grid,
      required this.clues,
      this.metadata,
      @JsonKey(name: 'created_at') required this.createdAt,
      @JsonKey(name: 'published_at') this.publishedAt})
      : _grid = grid,
        super._();

  factory _$PuzzleImpl.fromJson(Map<String, dynamic> json) =>
      _$$PuzzleImplFromJson(json);

  /// Unique identifier.
  @override
  final String id;

  /// Publication date (YYYY-MM-DD).
  @override
  final String date;

  /// Language code ("fr" or "en").
  @override
  final String language;

  /// Puzzle title.
  @override
  final String title;

  /// Author name.
  @override
  final String author;

  /// Difficulty level (1-5).
  @override
  final int difficulty;

  /// Status ("draft", "published", "archived").
  @override
  final String status;

  /// 2D grid of cells.
  final List<List<Cell>> _grid;

  /// 2D grid of cells.
  @override
  List<List<Cell>> get grid {
    if (_grid is EqualUnmodifiableListView) return _grid;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableListView(_grid);
  }

  /// Across and down clues.
  @override
  final Clues clues;

  /// Optional metadata.
  @override
  final Metadata? metadata;

  /// Creation timestamp.
  @override
  @JsonKey(name: 'created_at')
  final DateTime createdAt;

  /// Publication timestamp.
  @override
  @JsonKey(name: 'published_at')
  final DateTime? publishedAt;

  @override
  String toString() {
    return 'Puzzle(id: $id, date: $date, language: $language, title: $title, author: $author, difficulty: $difficulty, status: $status, grid: $grid, clues: $clues, metadata: $metadata, createdAt: $createdAt, publishedAt: $publishedAt)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$PuzzleImpl &&
            (identical(other.id, id) || other.id == id) &&
            (identical(other.date, date) || other.date == date) &&
            (identical(other.language, language) ||
                other.language == language) &&
            (identical(other.title, title) || other.title == title) &&
            (identical(other.author, author) || other.author == author) &&
            (identical(other.difficulty, difficulty) ||
                other.difficulty == difficulty) &&
            (identical(other.status, status) || other.status == status) &&
            const DeepCollectionEquality().equals(other._grid, _grid) &&
            (identical(other.clues, clues) || other.clues == clues) &&
            (identical(other.metadata, metadata) ||
                other.metadata == metadata) &&
            (identical(other.createdAt, createdAt) ||
                other.createdAt == createdAt) &&
            (identical(other.publishedAt, publishedAt) ||
                other.publishedAt == publishedAt));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(
      runtimeType,
      id,
      date,
      language,
      title,
      author,
      difficulty,
      status,
      const DeepCollectionEquality().hash(_grid),
      clues,
      metadata,
      createdAt,
      publishedAt);

  /// Create a copy of Puzzle
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$PuzzleImplCopyWith<_$PuzzleImpl> get copyWith =>
      __$$PuzzleImplCopyWithImpl<_$PuzzleImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$PuzzleImplToJson(
      this,
    );
  }
}

abstract class _Puzzle extends Puzzle {
  const factory _Puzzle(
          {required final String id,
          required final String date,
          required final String language,
          required final String title,
          required final String author,
          required final int difficulty,
          required final String status,
          required final List<List<Cell>> grid,
          required final Clues clues,
          final Metadata? metadata,
          @JsonKey(name: 'created_at') required final DateTime createdAt,
          @JsonKey(name: 'published_at') final DateTime? publishedAt}) =
      _$PuzzleImpl;
  const _Puzzle._() : super._();

  factory _Puzzle.fromJson(Map<String, dynamic> json) = _$PuzzleImpl.fromJson;

  /// Unique identifier.
  @override
  String get id;

  /// Publication date (YYYY-MM-DD).
  @override
  String get date;

  /// Language code ("fr" or "en").
  @override
  String get language;

  /// Puzzle title.
  @override
  String get title;

  /// Author name.
  @override
  String get author;

  /// Difficulty level (1-5).
  @override
  int get difficulty;

  /// Status ("draft", "published", "archived").
  @override
  String get status;

  /// 2D grid of cells.
  @override
  List<List<Cell>> get grid;

  /// Across and down clues.
  @override
  Clues get clues;

  /// Optional metadata.
  @override
  Metadata? get metadata;

  /// Creation timestamp.
  @override
  @JsonKey(name: 'created_at')
  DateTime get createdAt;

  /// Publication timestamp.
  @override
  @JsonKey(name: 'published_at')
  DateTime? get publishedAt;

  /// Create a copy of Puzzle
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$PuzzleImplCopyWith<_$PuzzleImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

PuzzleSummary _$PuzzleSummaryFromJson(Map<String, dynamic> json) {
  return _PuzzleSummary.fromJson(json);
}

/// @nodoc
mixin _$PuzzleSummary {
  String get id => throw _privateConstructorUsedError;
  String get date => throw _privateConstructorUsedError;
  String get language => throw _privateConstructorUsedError;
  String get title => throw _privateConstructorUsedError;
  String get author => throw _privateConstructorUsedError;
  int get difficulty => throw _privateConstructorUsedError;
  String get status => throw _privateConstructorUsedError;

  /// Serializes this PuzzleSummary to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of PuzzleSummary
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $PuzzleSummaryCopyWith<PuzzleSummary> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $PuzzleSummaryCopyWith<$Res> {
  factory $PuzzleSummaryCopyWith(
          PuzzleSummary value, $Res Function(PuzzleSummary) then) =
      _$PuzzleSummaryCopyWithImpl<$Res, PuzzleSummary>;
  @useResult
  $Res call(
      {String id,
      String date,
      String language,
      String title,
      String author,
      int difficulty,
      String status});
}

/// @nodoc
class _$PuzzleSummaryCopyWithImpl<$Res, $Val extends PuzzleSummary>
    implements $PuzzleSummaryCopyWith<$Res> {
  _$PuzzleSummaryCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of PuzzleSummary
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? date = null,
    Object? language = null,
    Object? title = null,
    Object? author = null,
    Object? difficulty = null,
    Object? status = null,
  }) {
    return _then(_value.copyWith(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as String,
      date: null == date
          ? _value.date
          : date // ignore: cast_nullable_to_non_nullable
              as String,
      language: null == language
          ? _value.language
          : language // ignore: cast_nullable_to_non_nullable
              as String,
      title: null == title
          ? _value.title
          : title // ignore: cast_nullable_to_non_nullable
              as String,
      author: null == author
          ? _value.author
          : author // ignore: cast_nullable_to_non_nullable
              as String,
      difficulty: null == difficulty
          ? _value.difficulty
          : difficulty // ignore: cast_nullable_to_non_nullable
              as int,
      status: null == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as String,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$PuzzleSummaryImplCopyWith<$Res>
    implements $PuzzleSummaryCopyWith<$Res> {
  factory _$$PuzzleSummaryImplCopyWith(
          _$PuzzleSummaryImpl value, $Res Function(_$PuzzleSummaryImpl) then) =
      __$$PuzzleSummaryImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {String id,
      String date,
      String language,
      String title,
      String author,
      int difficulty,
      String status});
}

/// @nodoc
class __$$PuzzleSummaryImplCopyWithImpl<$Res>
    extends _$PuzzleSummaryCopyWithImpl<$Res, _$PuzzleSummaryImpl>
    implements _$$PuzzleSummaryImplCopyWith<$Res> {
  __$$PuzzleSummaryImplCopyWithImpl(
      _$PuzzleSummaryImpl _value, $Res Function(_$PuzzleSummaryImpl) _then)
      : super(_value, _then);

  /// Create a copy of PuzzleSummary
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? date = null,
    Object? language = null,
    Object? title = null,
    Object? author = null,
    Object? difficulty = null,
    Object? status = null,
  }) {
    return _then(_$PuzzleSummaryImpl(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as String,
      date: null == date
          ? _value.date
          : date // ignore: cast_nullable_to_non_nullable
              as String,
      language: null == language
          ? _value.language
          : language // ignore: cast_nullable_to_non_nullable
              as String,
      title: null == title
          ? _value.title
          : title // ignore: cast_nullable_to_non_nullable
              as String,
      author: null == author
          ? _value.author
          : author // ignore: cast_nullable_to_non_nullable
              as String,
      difficulty: null == difficulty
          ? _value.difficulty
          : difficulty // ignore: cast_nullable_to_non_nullable
              as int,
      status: null == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as String,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$PuzzleSummaryImpl implements _PuzzleSummary {
  const _$PuzzleSummaryImpl(
      {required this.id,
      required this.date,
      required this.language,
      required this.title,
      required this.author,
      required this.difficulty,
      required this.status});

  factory _$PuzzleSummaryImpl.fromJson(Map<String, dynamic> json) =>
      _$$PuzzleSummaryImplFromJson(json);

  @override
  final String id;
  @override
  final String date;
  @override
  final String language;
  @override
  final String title;
  @override
  final String author;
  @override
  final int difficulty;
  @override
  final String status;

  @override
  String toString() {
    return 'PuzzleSummary(id: $id, date: $date, language: $language, title: $title, author: $author, difficulty: $difficulty, status: $status)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$PuzzleSummaryImpl &&
            (identical(other.id, id) || other.id == id) &&
            (identical(other.date, date) || other.date == date) &&
            (identical(other.language, language) ||
                other.language == language) &&
            (identical(other.title, title) || other.title == title) &&
            (identical(other.author, author) || other.author == author) &&
            (identical(other.difficulty, difficulty) ||
                other.difficulty == difficulty) &&
            (identical(other.status, status) || other.status == status));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(
      runtimeType, id, date, language, title, author, difficulty, status);

  /// Create a copy of PuzzleSummary
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$PuzzleSummaryImplCopyWith<_$PuzzleSummaryImpl> get copyWith =>
      __$$PuzzleSummaryImplCopyWithImpl<_$PuzzleSummaryImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$PuzzleSummaryImplToJson(
      this,
    );
  }
}

abstract class _PuzzleSummary implements PuzzleSummary {
  const factory _PuzzleSummary(
      {required final String id,
      required final String date,
      required final String language,
      required final String title,
      required final String author,
      required final int difficulty,
      required final String status}) = _$PuzzleSummaryImpl;

  factory _PuzzleSummary.fromJson(Map<String, dynamic> json) =
      _$PuzzleSummaryImpl.fromJson;

  @override
  String get id;
  @override
  String get date;
  @override
  String get language;
  @override
  String get title;
  @override
  String get author;
  @override
  int get difficulty;
  @override
  String get status;

  /// Create a copy of PuzzleSummary
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$PuzzleSummaryImplCopyWith<_$PuzzleSummaryImpl> get copyWith =>
      throw _privateConstructorUsedError;
}
