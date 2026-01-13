import 'package:freezed_annotation/freezed_annotation.dart';

part 'cell_model.freezed.dart';
part 'cell_model.g.dart';

/// Represents a cell in the crossword grid.
/// Supports both "mots croisés" (traditional) and "mots fléchés" (arrow words) formats.
/// Split clue cells can have both clueAcross and clueDown for two definitions.
@freezed
class Cell with _$Cell {
  const Cell._();

  const factory Cell({
    /// Cell type: "letter", "block", or "clue".
    required String type,

    /// Solution letter (A-Z) for letter cells.
    String? solution,

    /// Clue number if this cell starts an entry (mots croisés).
    int? number,

    /// Definition for across direction (→).
    @JsonKey(name: 'clue_across') String? clueAcross,

    /// Definition for down direction (↓).
    @JsonKey(name: 'clue_down') String? clueDown,
  }) = _Cell;

  factory Cell.fromJson(Map<String, dynamic> json) => _$CellFromJson(json);

  /// Whether this is a block cell.
  bool get isBlock => type == 'block';

  /// Whether this is a letter cell.
  bool get isLetter => type == 'letter';

  /// Whether this is a clue cell (mots fléchés).
  bool get isClue => type == 'clue';

  /// Whether this cell has a clue number.
  bool get hasNumber => number != null && number! > 0;

  /// Whether this is a split cell with both across and down clues.
  bool get isSplitClue => clueAcross != null && clueDown != null;

  /// Whether this cell has an across clue.
  bool get hasClueAcross => clueAcross != null;

  /// Whether this cell has a down clue.
  bool get hasClueDown => clueDown != null;
}
