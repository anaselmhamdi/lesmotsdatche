import 'package:freezed_annotation/freezed_annotation.dart';

part 'position_model.freezed.dart';
part 'position_model.g.dart';

/// Represents a position in the crossword grid.
@freezed
class Position with _$Position {
  const factory Position({
    required int row,
    required int col,
  }) = _Position;

  factory Position.fromJson(Map<String, dynamic> json) =>
      _$PositionFromJson(json);
}
