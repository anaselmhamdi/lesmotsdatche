/// API endpoint constants.
class ApiEndpoints {
  ApiEndpoints._();

  /// Health check endpoint.
  static const health = '/health';

  /// Get daily puzzle.
  static const dailyPuzzle = '/v1/puzzles/daily';

  /// Get puzzle by ID.
  static String puzzle(String id) => '/v1/puzzles/$id';

  /// List puzzles.
  static const puzzles = '/v1/puzzles';
}
