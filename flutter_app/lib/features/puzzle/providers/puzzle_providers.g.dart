// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'puzzle_providers.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

String _$puzzleRepositoryHash() => r'52962d9611ef980de549094ea9b76505decfdb03';

/// Provides the puzzle repository.
///
/// Copied from [puzzleRepository].
@ProviderFor(puzzleRepository)
final puzzleRepositoryProvider = AutoDisposeProvider<PuzzleRepository>.internal(
  puzzleRepository,
  name: r'puzzleRepositoryProvider',
  debugGetCreateSourceHash: const bool.fromEnvironment('dart.vm.product')
      ? null
      : _$puzzleRepositoryHash,
  dependencies: null,
  allTransitiveDependencies: null,
);

@Deprecated('Will be removed in 3.0. Use Ref instead')
// ignore: unused_element
typedef PuzzleRepositoryRef = AutoDisposeProviderRef<PuzzleRepository>;
String _$dailyPuzzleHash() => r'b308fd1eaaa704c473228ee5daba20bc7f052f5c';

/// Fetches the daily puzzle.
///
/// Copied from [dailyPuzzle].
@ProviderFor(dailyPuzzle)
final dailyPuzzleProvider = AutoDisposeFutureProvider<Puzzle>.internal(
  dailyPuzzle,
  name: r'dailyPuzzleProvider',
  debugGetCreateSourceHash:
      const bool.fromEnvironment('dart.vm.product') ? null : _$dailyPuzzleHash,
  dependencies: null,
  allTransitiveDependencies: null,
);

@Deprecated('Will be removed in 3.0. Use Ref instead')
// ignore: unused_element
typedef DailyPuzzleRef = AutoDisposeFutureProviderRef<Puzzle>;
String _$puzzleByIdHash() => r'7c2c44cfdfbcf5589d0fc0570da4f7b20ec619b9';

/// Copied from Dart SDK
class _SystemHash {
  _SystemHash._();

  static int combine(int hash, int value) {
    // ignore: parameter_assignments
    hash = 0x1fffffff & (hash + value);
    // ignore: parameter_assignments
    hash = 0x1fffffff & (hash + ((0x0007ffff & hash) << 10));
    return hash ^ (hash >> 6);
  }

  static int finish(int hash) {
    // ignore: parameter_assignments
    hash = 0x1fffffff & (hash + ((0x03ffffff & hash) << 3));
    // ignore: parameter_assignments
    hash = hash ^ (hash >> 11);
    return 0x1fffffff & (hash + ((0x00003fff & hash) << 15));
  }
}

/// Fetches a puzzle by ID.
///
/// Copied from [puzzleById].
@ProviderFor(puzzleById)
const puzzleByIdProvider = PuzzleByIdFamily();

/// Fetches a puzzle by ID.
///
/// Copied from [puzzleById].
class PuzzleByIdFamily extends Family<AsyncValue<Puzzle>> {
  /// Fetches a puzzle by ID.
  ///
  /// Copied from [puzzleById].
  const PuzzleByIdFamily();

  /// Fetches a puzzle by ID.
  ///
  /// Copied from [puzzleById].
  PuzzleByIdProvider call(
    String id,
  ) {
    return PuzzleByIdProvider(
      id,
    );
  }

  @override
  PuzzleByIdProvider getProviderOverride(
    covariant PuzzleByIdProvider provider,
  ) {
    return call(
      provider.id,
    );
  }

  static const Iterable<ProviderOrFamily>? _dependencies = null;

  @override
  Iterable<ProviderOrFamily>? get dependencies => _dependencies;

  static const Iterable<ProviderOrFamily>? _allTransitiveDependencies = null;

  @override
  Iterable<ProviderOrFamily>? get allTransitiveDependencies =>
      _allTransitiveDependencies;

  @override
  String? get name => r'puzzleByIdProvider';
}

/// Fetches a puzzle by ID.
///
/// Copied from [puzzleById].
class PuzzleByIdProvider extends AutoDisposeFutureProvider<Puzzle> {
  /// Fetches a puzzle by ID.
  ///
  /// Copied from [puzzleById].
  PuzzleByIdProvider(
    String id,
  ) : this._internal(
          (ref) => puzzleById(
            ref as PuzzleByIdRef,
            id,
          ),
          from: puzzleByIdProvider,
          name: r'puzzleByIdProvider',
          debugGetCreateSourceHash:
              const bool.fromEnvironment('dart.vm.product')
                  ? null
                  : _$puzzleByIdHash,
          dependencies: PuzzleByIdFamily._dependencies,
          allTransitiveDependencies:
              PuzzleByIdFamily._allTransitiveDependencies,
          id: id,
        );

  PuzzleByIdProvider._internal(
    super._createNotifier, {
    required super.name,
    required super.dependencies,
    required super.allTransitiveDependencies,
    required super.debugGetCreateSourceHash,
    required super.from,
    required this.id,
  }) : super.internal();

  final String id;

  @override
  Override overrideWith(
    FutureOr<Puzzle> Function(PuzzleByIdRef provider) create,
  ) {
    return ProviderOverride(
      origin: this,
      override: PuzzleByIdProvider._internal(
        (ref) => create(ref as PuzzleByIdRef),
        from: from,
        name: null,
        dependencies: null,
        allTransitiveDependencies: null,
        debugGetCreateSourceHash: null,
        id: id,
      ),
    );
  }

  @override
  AutoDisposeFutureProviderElement<Puzzle> createElement() {
    return _PuzzleByIdProviderElement(this);
  }

  @override
  bool operator ==(Object other) {
    return other is PuzzleByIdProvider && other.id == id;
  }

  @override
  int get hashCode {
    var hash = _SystemHash.combine(0, runtimeType.hashCode);
    hash = _SystemHash.combine(hash, id.hashCode);

    return _SystemHash.finish(hash);
  }
}

@Deprecated('Will be removed in 3.0. Use Ref instead')
// ignore: unused_element
mixin PuzzleByIdRef on AutoDisposeFutureProviderRef<Puzzle> {
  /// The parameter `id` of this provider.
  String get id;
}

class _PuzzleByIdProviderElement
    extends AutoDisposeFutureProviderElement<Puzzle> with PuzzleByIdRef {
  _PuzzleByIdProviderElement(super.provider);

  @override
  String get id => (origin as PuzzleByIdProvider).id;
}

String _$puzzleListHash() => r'9b08d1915a7b7790ff41857b854d1eff8c8ac15b';

/// Fetches the puzzle list.
///
/// Copied from [puzzleList].
@ProviderFor(puzzleList)
const puzzleListProvider = PuzzleListFamily();

/// Fetches the puzzle list.
///
/// Copied from [puzzleList].
class PuzzleListFamily extends Family<AsyncValue<List<PuzzleSummary>>> {
  /// Fetches the puzzle list.
  ///
  /// Copied from [puzzleList].
  const PuzzleListFamily();

  /// Fetches the puzzle list.
  ///
  /// Copied from [puzzleList].
  PuzzleListProvider call({
    String? language,
    int? limit,
  }) {
    return PuzzleListProvider(
      language: language,
      limit: limit,
    );
  }

  @override
  PuzzleListProvider getProviderOverride(
    covariant PuzzleListProvider provider,
  ) {
    return call(
      language: provider.language,
      limit: provider.limit,
    );
  }

  static const Iterable<ProviderOrFamily>? _dependencies = null;

  @override
  Iterable<ProviderOrFamily>? get dependencies => _dependencies;

  static const Iterable<ProviderOrFamily>? _allTransitiveDependencies = null;

  @override
  Iterable<ProviderOrFamily>? get allTransitiveDependencies =>
      _allTransitiveDependencies;

  @override
  String? get name => r'puzzleListProvider';
}

/// Fetches the puzzle list.
///
/// Copied from [puzzleList].
class PuzzleListProvider
    extends AutoDisposeFutureProvider<List<PuzzleSummary>> {
  /// Fetches the puzzle list.
  ///
  /// Copied from [puzzleList].
  PuzzleListProvider({
    String? language,
    int? limit,
  }) : this._internal(
          (ref) => puzzleList(
            ref as PuzzleListRef,
            language: language,
            limit: limit,
          ),
          from: puzzleListProvider,
          name: r'puzzleListProvider',
          debugGetCreateSourceHash:
              const bool.fromEnvironment('dart.vm.product')
                  ? null
                  : _$puzzleListHash,
          dependencies: PuzzleListFamily._dependencies,
          allTransitiveDependencies:
              PuzzleListFamily._allTransitiveDependencies,
          language: language,
          limit: limit,
        );

  PuzzleListProvider._internal(
    super._createNotifier, {
    required super.name,
    required super.dependencies,
    required super.allTransitiveDependencies,
    required super.debugGetCreateSourceHash,
    required super.from,
    required this.language,
    required this.limit,
  }) : super.internal();

  final String? language;
  final int? limit;

  @override
  Override overrideWith(
    FutureOr<List<PuzzleSummary>> Function(PuzzleListRef provider) create,
  ) {
    return ProviderOverride(
      origin: this,
      override: PuzzleListProvider._internal(
        (ref) => create(ref as PuzzleListRef),
        from: from,
        name: null,
        dependencies: null,
        allTransitiveDependencies: null,
        debugGetCreateSourceHash: null,
        language: language,
        limit: limit,
      ),
    );
  }

  @override
  AutoDisposeFutureProviderElement<List<PuzzleSummary>> createElement() {
    return _PuzzleListProviderElement(this);
  }

  @override
  bool operator ==(Object other) {
    return other is PuzzleListProvider &&
        other.language == language &&
        other.limit == limit;
  }

  @override
  int get hashCode {
    var hash = _SystemHash.combine(0, runtimeType.hashCode);
    hash = _SystemHash.combine(hash, language.hashCode);
    hash = _SystemHash.combine(hash, limit.hashCode);

    return _SystemHash.finish(hash);
  }
}

@Deprecated('Will be removed in 3.0. Use Ref instead')
// ignore: unused_element
mixin PuzzleListRef on AutoDisposeFutureProviderRef<List<PuzzleSummary>> {
  /// The parameter `language` of this provider.
  String? get language;

  /// The parameter `limit` of this provider.
  int? get limit;
}

class _PuzzleListProviderElement
    extends AutoDisposeFutureProviderElement<List<PuzzleSummary>>
    with PuzzleListRef {
  _PuzzleListProviderElement(super.provider);

  @override
  String? get language => (origin as PuzzleListProvider).language;
  @override
  int? get limit => (origin as PuzzleListProvider).limit;
}
// ignore_for_file: type=lint
// ignore_for_file: subtype_of_sealed_class, invalid_use_of_internal_member, invalid_use_of_visible_for_testing_member, deprecated_member_use_from_same_package
