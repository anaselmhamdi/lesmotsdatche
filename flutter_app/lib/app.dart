import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'features/puzzle/presentation/screens/home_screen.dart';
import 'features/puzzle/presentation/screens/puzzle_player_screen.dart';

/// Main application widget.
class LesMotsApp extends StatelessWidget {
  const LesMotsApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      title: "Les Mots d'Atché",
      debugShowCheckedModeBanner: false,
      theme: _buildLightTheme(),
      darkTheme: _buildDarkTheme(),
      themeMode: ThemeMode.system,
      routerConfig: _router,
      locale: const Locale('fr'),
      supportedLocales: const [
        Locale('fr'),
      ],
    );
  }

  ThemeData _buildLightTheme() {
    return ThemeData(
      useMaterial3: true,
      colorScheme: ColorScheme.fromSeed(
        seedColor: const Color(0xFF1E88E5), // Blue
        brightness: Brightness.light,
      ),
      appBarTheme: const AppBarTheme(
        centerTitle: true,
        elevation: 0,
      ),
      cardTheme: CardThemeData(
        elevation: 2,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
        ),
      ),
    );
  }

  ThemeData _buildDarkTheme() {
    return ThemeData(
      useMaterial3: true,
      colorScheme: ColorScheme.fromSeed(
        seedColor: const Color(0xFF1E88E5), // Blue
        brightness: Brightness.dark,
      ),
      appBarTheme: const AppBarTheme(
        centerTitle: true,
        elevation: 0,
      ),
      cardTheme: CardThemeData(
        elevation: 2,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
        ),
      ),
    );
  }
}

/// Application router configuration.
final _router = GoRouter(
  initialLocation: '/',
  routes: [
    GoRoute(
      path: '/',
      name: 'home',
      builder: (context, state) => const HomeScreen(),
    ),
    GoRoute(
      path: '/puzzle/:id',
      name: 'puzzle',
      builder: (context, state) => PuzzlePlayerScreen(
        puzzleId: state.pathParameters['id']!,
      ),
    ),
  ],
);
