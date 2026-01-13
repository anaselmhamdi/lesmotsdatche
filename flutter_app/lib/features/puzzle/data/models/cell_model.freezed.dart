// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'cell_model.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
    'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more information: https://github.com/rrousselGit/freezed#adding-getters-and-methods-to-our-models');

Cell _$CellFromJson(Map<String, dynamic> json) {
  return _Cell.fromJson(json);
}

/// @nodoc
mixin _$Cell {
  /// Cell type: "letter", "block", or "clue".
  String get type => throw _privateConstructorUsedError;

  /// Solution letter (A-Z) for letter cells.
  String? get solution => throw _privateConstructorUsedError;

  /// Clue number if this cell starts an entry (mots croisés).
  int? get number => throw _privateConstructorUsedError;

  /// Definition for across direction (→).
  @JsonKey(name: 'clue_across')
  String? get clueAcross => throw _privateConstructorUsedError;

  /// Definition for down direction (↓).
  @JsonKey(name: 'clue_down')
  String? get clueDown => throw _privateConstructorUsedError;

  /// Serializes this Cell to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of Cell
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $CellCopyWith<Cell> get copyWith => throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $CellCopyWith<$Res> {
  factory $CellCopyWith(Cell value, $Res Function(Cell) then) =
      _$CellCopyWithImpl<$Res, Cell>;
  @useResult
  $Res call(
      {String type,
      String? solution,
      int? number,
      @JsonKey(name: 'clue_across') String? clueAcross,
      @JsonKey(name: 'clue_down') String? clueDown});
}

/// @nodoc
class _$CellCopyWithImpl<$Res, $Val extends Cell>
    implements $CellCopyWith<$Res> {
  _$CellCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of Cell
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? type = null,
    Object? solution = freezed,
    Object? number = freezed,
    Object? clueAcross = freezed,
    Object? clueDown = freezed,
  }) {
    return _then(_value.copyWith(
      type: null == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as String,
      solution: freezed == solution
          ? _value.solution
          : solution // ignore: cast_nullable_to_non_nullable
              as String?,
      number: freezed == number
          ? _value.number
          : number // ignore: cast_nullable_to_non_nullable
              as int?,
      clueAcross: freezed == clueAcross
          ? _value.clueAcross
          : clueAcross // ignore: cast_nullable_to_non_nullable
              as String?,
      clueDown: freezed == clueDown
          ? _value.clueDown
          : clueDown // ignore: cast_nullable_to_non_nullable
              as String?,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$CellImplCopyWith<$Res> implements $CellCopyWith<$Res> {
  factory _$$CellImplCopyWith(
          _$CellImpl value, $Res Function(_$CellImpl) then) =
      __$$CellImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {String type,
      String? solution,
      int? number,
      @JsonKey(name: 'clue_across') String? clueAcross,
      @JsonKey(name: 'clue_down') String? clueDown});
}

/// @nodoc
class __$$CellImplCopyWithImpl<$Res>
    extends _$CellCopyWithImpl<$Res, _$CellImpl>
    implements _$$CellImplCopyWith<$Res> {
  __$$CellImplCopyWithImpl(_$CellImpl _value, $Res Function(_$CellImpl) _then)
      : super(_value, _then);

  /// Create a copy of Cell
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? type = null,
    Object? solution = freezed,
    Object? number = freezed,
    Object? clueAcross = freezed,
    Object? clueDown = freezed,
  }) {
    return _then(_$CellImpl(
      type: null == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as String,
      solution: freezed == solution
          ? _value.solution
          : solution // ignore: cast_nullable_to_non_nullable
              as String?,
      number: freezed == number
          ? _value.number
          : number // ignore: cast_nullable_to_non_nullable
              as int?,
      clueAcross: freezed == clueAcross
          ? _value.clueAcross
          : clueAcross // ignore: cast_nullable_to_non_nullable
              as String?,
      clueDown: freezed == clueDown
          ? _value.clueDown
          : clueDown // ignore: cast_nullable_to_non_nullable
              as String?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$CellImpl extends _Cell {
  const _$CellImpl(
      {required this.type,
      this.solution,
      this.number,
      @JsonKey(name: 'clue_across') this.clueAcross,
      @JsonKey(name: 'clue_down') this.clueDown})
      : super._();

  factory _$CellImpl.fromJson(Map<String, dynamic> json) =>
      _$$CellImplFromJson(json);

  /// Cell type: "letter", "block", or "clue".
  @override
  final String type;

  /// Solution letter (A-Z) for letter cells.
  @override
  final String? solution;

  /// Clue number if this cell starts an entry (mots croisés).
  @override
  final int? number;

  /// Definition for across direction (→).
  @override
  @JsonKey(name: 'clue_across')
  final String? clueAcross;

  /// Definition for down direction (↓).
  @override
  @JsonKey(name: 'clue_down')
  final String? clueDown;

  @override
  String toString() {
    return 'Cell(type: $type, solution: $solution, number: $number, clueAcross: $clueAcross, clueDown: $clueDown)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$CellImpl &&
            (identical(other.type, type) || other.type == type) &&
            (identical(other.solution, solution) ||
                other.solution == solution) &&
            (identical(other.number, number) || other.number == number) &&
            (identical(other.clueAcross, clueAcross) ||
                other.clueAcross == clueAcross) &&
            (identical(other.clueDown, clueDown) ||
                other.clueDown == clueDown));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode =>
      Object.hash(runtimeType, type, solution, number, clueAcross, clueDown);

  /// Create a copy of Cell
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$CellImplCopyWith<_$CellImpl> get copyWith =>
      __$$CellImplCopyWithImpl<_$CellImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$CellImplToJson(
      this,
    );
  }
}

abstract class _Cell extends Cell {
  const factory _Cell(
      {required final String type,
      final String? solution,
      final int? number,
      @JsonKey(name: 'clue_across') final String? clueAcross,
      @JsonKey(name: 'clue_down') final String? clueDown}) = _$CellImpl;
  const _Cell._() : super._();

  factory _Cell.fromJson(Map<String, dynamic> json) = _$CellImpl.fromJson;

  /// Cell type: "letter", "block", or "clue".
  @override
  String get type;

  /// Solution letter (A-Z) for letter cells.
  @override
  String? get solution;

  /// Clue number if this cell starts an entry (mots croisés).
  @override
  int? get number;

  /// Definition for across direction (→).
  @override
  @JsonKey(name: 'clue_across')
  String? get clueAcross;

  /// Definition for down direction (↓).
  @override
  @JsonKey(name: 'clue_down')
  String? get clueDown;

  /// Create a copy of Cell
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$CellImplCopyWith<_$CellImpl> get copyWith =>
      throw _privateConstructorUsedError;
}
