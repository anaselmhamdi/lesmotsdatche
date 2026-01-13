import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:lesmotsdatche/core/api/api_client.dart';
import 'package:lesmotsdatche/features/puzzle/data/repositories/puzzle_repository.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

@GenerateMocks([ApiClient])
import 'puzzle_repository_test.mocks.dart';

void main() {
  late MockApiClient mockClient;
  late PuzzleRepository repository;

  setUp(() {
    mockClient = MockApiClient();
    repository = PuzzleRepository(mockClient);
  });

  group('PuzzleRepository', () {
    group('getDailyPuzzle', () {
      test('returns puzzle when successful', () async {
        final puzzleJson = {
          'id': 'daily-puzzle',
          'date': '2024-01-15',
          'language': 'fr',
          'title': 'Daily Puzzle',
          'author': 'Test',
          'difficulty': 3,
          'status': 'published',
          'grid': [
            [
              {'type': 'letter', 'solution': 'A'},
            ],
          ],
          'clues': {
            'across': [],
            'down': [],
          },
          'created_at': '2024-01-15T10:00:00Z',
        };

        when(mockClient.get(
          '/v1/puzzles/daily',
          queryParams: {'language': 'fr'},
        )).thenAnswer((_) async => Response(
              data: puzzleJson,
              statusCode: 200,
              requestOptions: RequestOptions(),
            ));

        final puzzle = await repository.getDailyPuzzle();

        expect(puzzle.id, 'daily-puzzle');
        expect(puzzle.title, 'Daily Puzzle');
        verify(mockClient.get(
          '/v1/puzzles/daily',
          queryParams: {'language': 'fr'},
        )).called(1);
      });
    });

    group('getPuzzleById', () {
      test('returns puzzle when found', () async {
        final puzzleJson = {
          'id': 'test-id',
          'date': '2024-01-15',
          'language': 'fr',
          'title': 'Test Puzzle',
          'author': 'Test',
          'difficulty': 3,
          'status': 'published',
          'grid': [
            [
              {'type': 'letter', 'solution': 'A'},
            ],
          ],
          'clues': {
            'across': [],
            'down': [],
          },
          'created_at': '2024-01-15T10:00:00Z',
        };

        when(mockClient.get('/v1/puzzles/test-id'))
            .thenAnswer((_) async => Response(
                  data: puzzleJson,
                  statusCode: 200,
                  requestOptions: RequestOptions(),
                ));

        final puzzle = await repository.getPuzzleById('test-id');

        expect(puzzle.id, 'test-id');
        verify(mockClient.get('/v1/puzzles/test-id')).called(1);
      });
    });

    group('listPuzzles', () {
      test('returns list of puzzle summaries', () async {
        final responseJson = {
          'puzzles': [
            {
              'id': 'puzzle-1',
              'date': '2024-01-15',
              'language': 'fr',
              'title': 'Puzzle 1',
              'author': 'Author',
              'difficulty': 2,
              'status': 'published',
            },
            {
              'id': 'puzzle-2',
              'date': '2024-01-14',
              'language': 'fr',
              'title': 'Puzzle 2',
              'author': 'Author',
              'difficulty': 3,
              'status': 'published',
            },
          ],
          'count': 2,
        };

        when(mockClient.get(
          '/v1/puzzles',
          queryParams: anyNamed('queryParams'),
        )).thenAnswer((_) async => Response(
              data: responseJson,
              statusCode: 200,
              requestOptions: RequestOptions(),
            ));

        final summaries = await repository.listPuzzles(language: 'fr');

        expect(summaries.length, 2);
        expect(summaries[0].id, 'puzzle-1');
        expect(summaries[1].id, 'puzzle-2');
      });
    });
  });
}
