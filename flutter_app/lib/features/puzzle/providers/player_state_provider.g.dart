// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'player_state_provider.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

String _$playerStateNotifierHash() =>
    r'41b87162808551b86a502e9edba84f0724d533a3';

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

abstract class _$PlayerStateNotifier
    extends BuildlessAutoDisposeNotifier<PlayerState> {
  late final Puzzle puzzle;

  PlayerState build(
    Puzzle puzzle,
  );
}

/// Manages the player state for a puzzle.
///
/// Copied from [PlayerStateNotifier].
@ProviderFor(PlayerStateNotifier)
const playerStateNotifierProvider = PlayerStateNotifierFamily();

/// Manages the player state for a puzzle.
///
/// Copied from [PlayerStateNotifier].
class PlayerStateNotifierFamily extends Family<PlayerState> {
  /// Manages the player state for a puzzle.
  ///
  /// Copied from [PlayerStateNotifier].
  const PlayerStateNotifierFamily();

  /// Manages the player state for a puzzle.
  ///
  /// Copied from [PlayerStateNotifier].
  PlayerStateNotifierProvider call(
    Puzzle puzzle,
  ) {
    return PlayerStateNotifierProvider(
      puzzle,
    );
  }

  @override
  PlayerStateNotifierProvider getProviderOverride(
    covariant PlayerStateNotifierProvider provider,
  ) {
    return call(
      provider.puzzle,
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
  String? get name => r'playerStateNotifierProvider';
}

/// Manages the player state for a puzzle.
///
/// Copied from [PlayerStateNotifier].
class PlayerStateNotifierProvider
    extends AutoDisposeNotifierProviderImpl<PlayerStateNotifier, PlayerState> {
  /// Manages the player state for a puzzle.
  ///
  /// Copied from [PlayerStateNotifier].
  PlayerStateNotifierProvider(
    Puzzle puzzle,
  ) : this._internal(
          () => PlayerStateNotifier()..puzzle = puzzle,
          from: playerStateNotifierProvider,
          name: r'playerStateNotifierProvider',
          debugGetCreateSourceHash:
              const bool.fromEnvironment('dart.vm.product')
                  ? null
                  : _$playerStateNotifierHash,
          dependencies: PlayerStateNotifierFamily._dependencies,
          allTransitiveDependencies:
              PlayerStateNotifierFamily._allTransitiveDependencies,
          puzzle: puzzle,
        );

  PlayerStateNotifierProvider._internal(
    super._createNotifier, {
    required super.name,
    required super.dependencies,
    required super.allTransitiveDependencies,
    required super.debugGetCreateSourceHash,
    required super.from,
    required this.puzzle,
  }) : super.internal();

  final Puzzle puzzle;

  @override
  PlayerState runNotifierBuild(
    covariant PlayerStateNotifier notifier,
  ) {
    return notifier.build(
      puzzle,
    );
  }

  @override
  Override overrideWith(PlayerStateNotifier Function() create) {
    return ProviderOverride(
      origin: this,
      override: PlayerStateNotifierProvider._internal(
        () => create()..puzzle = puzzle,
        from: from,
        name: null,
        dependencies: null,
        allTransitiveDependencies: null,
        debugGetCreateSourceHash: null,
        puzzle: puzzle,
      ),
    );
  }

  @override
  AutoDisposeNotifierProviderElement<PlayerStateNotifier, PlayerState>
      createElement() {
    return _PlayerStateNotifierProviderElement(this);
  }

  @override
  bool operator ==(Object other) {
    return other is PlayerStateNotifierProvider && other.puzzle == puzzle;
  }

  @override
  int get hashCode {
    var hash = _SystemHash.combine(0, runtimeType.hashCode);
    hash = _SystemHash.combine(hash, puzzle.hashCode);

    return _SystemHash.finish(hash);
  }
}

@Deprecated('Will be removed in 3.0. Use Ref instead')
// ignore: unused_element
mixin PlayerStateNotifierRef on AutoDisposeNotifierProviderRef<PlayerState> {
  /// The parameter `puzzle` of this provider.
  Puzzle get puzzle;
}

class _PlayerStateNotifierProviderElement
    extends AutoDisposeNotifierProviderElement<PlayerStateNotifier, PlayerState>
    with PlayerStateNotifierRef {
  _PlayerStateNotifierProviderElement(super.provider);

  @override
  Puzzle get puzzle => (origin as PlayerStateNotifierProvider).puzzle;
}
// ignore_for_file: type=lint
// ignore_for_file: subtype_of_sealed_class, invalid_use_of_internal_member, invalid_use_of_visible_for_testing_member, deprecated_member_use_from_same_package
