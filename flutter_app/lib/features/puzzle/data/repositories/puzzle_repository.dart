import '../../../../core/api/api_client.dart';
import '../../../../core/api/api_endpoints.dart';
import '../models/models.dart';

/// Repository for puzzle-related API operations.
class PuzzleRepository {
  final ApiClient _client;

  PuzzleRepository(this._client);

  /// Fetches the daily puzzle.
  Future<Puzzle> getDailyPuzzle({String language = 'fr'}) async {
    final response = await _client.get(
      ApiEndpoints.dailyPuzzle,
      queryParams: {'language': language},
    );
    return Puzzle.fromJson(response.data as Map<String, dynamic>);
  }

  /// Fetches a puzzle by ID.
  Future<Puzzle> getPuzzleById(String id) async {
    final response = await _client.get(ApiEndpoints.puzzle(id));
    return Puzzle.fromJson(response.data as Map<String, dynamic>);
  }

  /// Lists puzzles with optional filters.
  Future<List<PuzzleSummary>> listPuzzles({
    String? language,
    String? from,
    String? to,
    int? difficulty,
    int? limit,
  }) async {
    final response = await _client.get(
      ApiEndpoints.puzzles,
      queryParams: {
        if (language != null) 'language': language,
        if (from != null) 'from': from,
        if (to != null) 'to': to,
        if (difficulty != null) 'difficulty': difficulty,
        if (limit != null) 'limit': limit,
      },
    );

    final data = response.data as Map<String, dynamic>;
    final puzzles = data['puzzles'] as List<dynamic>;
    return puzzles
        .map((e) => PuzzleSummary.fromJson(e as Map<String, dynamic>))
        .toList();
  }
}
