import 'package:flutter_test/flutter_test.dart';
import 'package:lesmotsdatche/features/puzzle/data/models/models.dart';

void main() {
  group('Position', () {
    test('fromJson creates Position correctly', () {
      final json = {'row': 1, 'col': 2};
      final position = Position.fromJson(json);

      expect(position.row, 1);
      expect(position.col, 2);
    });

    test('toJson serializes correctly', () {
      const position = Position(row: 3, col: 4);
      final json = position.toJson();

      expect(json['row'], 3);
      expect(json['col'], 4);
    });
  });

  group('Cell', () {
    test('fromJson creates letter cell correctly', () {
      final json = {
        'type': 'letter',
        'solution': 'A',
        'number': 1,
      };
      final cell = Cell.fromJson(json);

      expect(cell.type, 'letter');
      expect(cell.solution, 'A');
      expect(cell.number, 1);
      expect(cell.isLetter, true);
      expect(cell.isBlock, false);
    });

    test('fromJson creates block cell correctly', () {
      final json = {'type': 'block'};
      final cell = Cell.fromJson(json);

      expect(cell.type, 'block');
      expect(cell.solution, isNull);
      expect(cell.isBlock, true);
      expect(cell.isLetter, false);
    });
  });

  group('Clue', () {
    test('fromJson creates clue correctly', () {
      final json = {
        'id': '1-across',
        'direction': 'across',
        'number': 1,
        'prompt': 'Test clue',
        'answer': 'TEST',
        'start': {'row': 0, 'col': 0},
        'length': 4,
      };
      final clue = Clue.fromJson(json);

      expect(clue.id, '1-across');
      expect(clue.direction, 'across');
      expect(clue.number, 1);
      expect(clue.prompt, 'Test clue');
      expect(clue.answer, 'TEST');
      expect(clue.start.row, 0);
      expect(clue.length, 4);
      expect(clue.isAcross, true);
      expect(clue.isDown, false);
    });
  });

  group('Puzzle', () {
    test('fromJson creates puzzle correctly', () {
      final json = {
        'id': 'test-puzzle',
        'date': '2024-01-15',
        'language': 'fr',
        'title': 'Test Puzzle',
        'author': 'Test Author',
        'difficulty': 3,
        'status': 'published',
        'grid': [
          [
            {'type': 'letter', 'solution': 'A', 'number': 1},
            {'type': 'letter', 'solution': 'B'},
          ],
          [
            {'type': 'letter', 'solution': 'C'},
            {'type': 'block'},
          ],
        ],
        'clues': {
          'across': [
            {
              'id': '1-across',
              'direction': 'across',
              'number': 1,
              'prompt': 'Test',
              'answer': 'AB',
              'start': {'row': 0, 'col': 0},
              'length': 2,
            },
          ],
          'down': [],
        },
        'created_at': '2024-01-15T10:00:00Z',
      };

      final puzzle = Puzzle.fromJson(json);

      expect(puzzle.id, 'test-puzzle');
      expect(puzzle.date, '2024-01-15');
      expect(puzzle.language, 'fr');
      expect(puzzle.title, 'Test Puzzle');
      expect(puzzle.author, 'Test Author');
      expect(puzzle.difficulty, 3);
      expect(puzzle.rows, 2);
      expect(puzzle.cols, 2);
      expect(puzzle.totalClues, 1);
      expect(puzzle.isPublished, true);
    });
  });

  group('PuzzleSummary', () {
    test('fromJson creates summary correctly', () {
      final json = {
        'id': 'test-id',
        'date': '2024-01-15',
        'language': 'fr',
        'title': 'Test',
        'author': 'Author',
        'difficulty': 2,
        'status': 'published',
      };

      final summary = PuzzleSummary.fromJson(json);

      expect(summary.id, 'test-id');
      expect(summary.title, 'Test');
      expect(summary.difficulty, 2);
    });
  });
}
