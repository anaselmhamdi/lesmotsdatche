import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../../../core/providers/core_providers.dart';
import '../data/models/models.dart';
import '../data/repositories/puzzle_repository.dart';

part 'puzzle_providers.g.dart';

/// Provides the puzzle repository.
@riverpod
PuzzleRepository puzzleRepository(PuzzleRepositoryRef ref) {
  final client = ref.watch(apiClientProvider);
  return PuzzleRepository(client);
}

/// Fetches the daily puzzle.
@riverpod
Future<Puzzle> dailyPuzzle(DailyPuzzleRef ref) async {
  final repo = ref.watch(puzzleRepositoryProvider);
  return repo.getDailyPuzzle();
}

/// Fetches a puzzle by ID.
@riverpod
Future<Puzzle> puzzleById(PuzzleByIdRef ref, String id) async {
  final repo = ref.watch(puzzleRepositoryProvider);
  return repo.getPuzzleById(id);
}

/// Fetches the puzzle list.
@riverpod
Future<List<PuzzleSummary>> puzzleList(
  PuzzleListRef ref, {
  String? language,
  int? limit,
}) async {
  final repo = ref.watch(puzzleRepositoryProvider);
  return repo.listPuzzles(language: language, limit: limit);
}
