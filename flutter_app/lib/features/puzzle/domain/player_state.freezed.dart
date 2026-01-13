// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'player_state.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
    'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more information: https://github.com/rrousselGit/freezed#adding-getters-and-methods-to-our-models');

/// @nodoc
mixin _$PlayerState {
  /// The puzzle being played.
  Puzzle get puzzle => throw _privateConstructorUsedError;

  /// User's letter entries (mirrors grid dimensions).
  List<List<String>> get userInput => throw _privateConstructorUsedError;

  /// Currently selected cell position.
  Position? get selectedCell => throw _privateConstructorUsedError;

  /// Current navigation direction.
  ClueDirection get direction => throw _privateConstructorUsedError;

  /// Cells that have been correctly filled.
  Set<Position> get completedCells => throw _privateConstructorUsedError;

  /// Create a copy of PlayerState
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $PlayerStateCopyWith<PlayerState> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $PlayerStateCopyWith<$Res> {
  factory $PlayerStateCopyWith(
          PlayerState value, $Res Function(PlayerState) then) =
      _$PlayerStateCopyWithImpl<$Res, PlayerState>;
  @useResult
  $Res call(
      {Puzzle puzzle,
      List<List<String>> userInput,
      Position? selectedCell,
      ClueDirection direction,
      Set<Position> completedCells});

  $PuzzleCopyWith<$Res> get puzzle;
  $PositionCopyWith<$Res>? get selectedCell;
}

/// @nodoc
class _$PlayerStateCopyWithImpl<$Res, $Val extends PlayerState>
    implements $PlayerStateCopyWith<$Res> {
  _$PlayerStateCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of PlayerState
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? puzzle = null,
    Object? userInput = null,
    Object? selectedCell = freezed,
    Object? direction = null,
    Object? completedCells = null,
  }) {
    return _then(_value.copyWith(
      puzzle: null == puzzle
          ? _value.puzzle
          : puzzle // ignore: cast_nullable_to_non_nullable
              as Puzzle,
      userInput: null == userInput
          ? _value.userInput
          : userInput // ignore: cast_nullable_to_non_nullable
              as List<List<String>>,
      selectedCell: freezed == selectedCell
          ? _value.selectedCell
          : selectedCell // ignore: cast_nullable_to_non_nullable
              as Position?,
      direction: null == direction
          ? _value.direction
          : direction // ignore: cast_nullable_to_non_nullable
              as ClueDirection,
      completedCells: null == completedCells
          ? _value.completedCells
          : completedCells // ignore: cast_nullable_to_non_nullable
              as Set<Position>,
    ) as $Val);
  }

  /// Create a copy of PlayerState
  /// with the given fields replaced by the non-null parameter values.
  @override
  @pragma('vm:prefer-inline')
  $PuzzleCopyWith<$Res> get puzzle {
    return $PuzzleCopyWith<$Res>(_value.puzzle, (value) {
      return _then(_value.copyWith(puzzle: value) as $Val);
    });
  }

  /// Create a copy of PlayerState
  /// with the given fields replaced by the non-null parameter values.
  @override
  @pragma('vm:prefer-inline')
  $PositionCopyWith<$Res>? get selectedCell {
    if (_value.selectedCell == null) {
      return null;
    }

    return $PositionCopyWith<$Res>(_value.selectedCell!, (value) {
      return _then(_value.copyWith(selectedCell: value) as $Val);
    });
  }
}

/// @nodoc
abstract class _$$PlayerStateImplCopyWith<$Res>
    implements $PlayerStateCopyWith<$Res> {
  factory _$$PlayerStateImplCopyWith(
          _$PlayerStateImpl value, $Res Function(_$PlayerStateImpl) then) =
      __$$PlayerStateImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {Puzzle puzzle,
      List<List<String>> userInput,
      Position? selectedCell,
      ClueDirection direction,
      Set<Position> completedCells});

  @override
  $PuzzleCopyWith<$Res> get puzzle;
  @override
  $PositionCopyWith<$Res>? get selectedCell;
}

/// @nodoc
class __$$PlayerStateImplCopyWithImpl<$Res>
    extends _$PlayerStateCopyWithImpl<$Res, _$PlayerStateImpl>
    implements _$$PlayerStateImplCopyWith<$Res> {
  __$$PlayerStateImplCopyWithImpl(
      _$PlayerStateImpl _value, $Res Function(_$PlayerStateImpl) _then)
      : super(_value, _then);

  /// Create a copy of PlayerState
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? puzzle = null,
    Object? userInput = null,
    Object? selectedCell = freezed,
    Object? direction = null,
    Object? completedCells = null,
  }) {
    return _then(_$PlayerStateImpl(
      puzzle: null == puzzle
          ? _value.puzzle
          : puzzle // ignore: cast_nullable_to_non_nullable
              as Puzzle,
      userInput: null == userInput
          ? _value._userInput
          : userInput // ignore: cast_nullable_to_non_nullable
              as List<List<String>>,
      selectedCell: freezed == selectedCell
          ? _value.selectedCell
          : selectedCell // ignore: cast_nullable_to_non_nullable
              as Position?,
      direction: null == direction
          ? _value.direction
          : direction // ignore: cast_nullable_to_non_nullable
              as ClueDirection,
      completedCells: null == completedCells
          ? _value._completedCells
          : completedCells // ignore: cast_nullable_to_non_nullable
              as Set<Position>,
    ));
  }
}

/// @nodoc

class _$PlayerStateImpl extends _PlayerState {
  const _$PlayerStateImpl(
      {required this.puzzle,
      required final List<List<String>> userInput,
      this.selectedCell,
      this.direction = ClueDirection.across,
      final Set<Position> completedCells = const {}})
      : _userInput = userInput,
        _completedCells = completedCells,
        super._();

  /// The puzzle being played.
  @override
  final Puzzle puzzle;

  /// User's letter entries (mirrors grid dimensions).
  final List<List<String>> _userInput;

  /// User's letter entries (mirrors grid dimensions).
  @override
  List<List<String>> get userInput {
    if (_userInput is EqualUnmodifiableListView) return _userInput;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableListView(_userInput);
  }

  /// Currently selected cell position.
  @override
  final Position? selectedCell;

  /// Current navigation direction.
  @override
  @JsonKey()
  final ClueDirection direction;

  /// Cells that have been correctly filled.
  final Set<Position> _completedCells;

  /// Cells that have been correctly filled.
  @override
  @JsonKey()
  Set<Position> get completedCells {
    if (_completedCells is EqualUnmodifiableSetView) return _completedCells;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableSetView(_completedCells);
  }

  @override
  String toString() {
    return 'PlayerState(puzzle: $puzzle, userInput: $userInput, selectedCell: $selectedCell, direction: $direction, completedCells: $completedCells)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$PlayerStateImpl &&
            (identical(other.puzzle, puzzle) || other.puzzle == puzzle) &&
            const DeepCollectionEquality()
                .equals(other._userInput, _userInput) &&
            (identical(other.selectedCell, selectedCell) ||
                other.selectedCell == selectedCell) &&
            (identical(other.direction, direction) ||
                other.direction == direction) &&
            const DeepCollectionEquality()
                .equals(other._completedCells, _completedCells));
  }

  @override
  int get hashCode => Object.hash(
      runtimeType,
      puzzle,
      const DeepCollectionEquality().hash(_userInput),
      selectedCell,
      direction,
      const DeepCollectionEquality().hash(_completedCells));

  /// Create a copy of PlayerState
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$PlayerStateImplCopyWith<_$PlayerStateImpl> get copyWith =>
      __$$PlayerStateImplCopyWithImpl<_$PlayerStateImpl>(this, _$identity);
}

abstract class _PlayerState extends PlayerState {
  const factory _PlayerState(
      {required final Puzzle puzzle,
      required final List<List<String>> userInput,
      final Position? selectedCell,
      final ClueDirection direction,
      final Set<Position> completedCells}) = _$PlayerStateImpl;
  const _PlayerState._() : super._();

  /// The puzzle being played.
  @override
  Puzzle get puzzle;

  /// User's letter entries (mirrors grid dimensions).
  @override
  List<List<String>> get userInput;

  /// Currently selected cell position.
  @override
  Position? get selectedCell;

  /// Current navigation direction.
  @override
  ClueDirection get direction;

  /// Cells that have been correctly filled.
  @override
  Set<Position> get completedCells;

  /// Create a copy of PlayerState
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$PlayerStateImplCopyWith<_$PlayerStateImpl> get copyWith =>
      throw _privateConstructorUsedError;
}
