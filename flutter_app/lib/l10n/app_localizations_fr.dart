// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for French (`fr`).
class AppLocalizationsFr extends AppLocalizations {
  AppLocalizationsFr([String locale = 'fr']) : super(locale);

  @override
  String get appTitle => 'Les Mots d\'Atché';

  @override
  String get dailyPuzzle => 'Puzzle du jour';

  @override
  String get loading => 'Chargement...';

  @override
  String get errorLoading => 'Erreur de chargement';

  @override
  String get retry => 'Réessayer';

  @override
  String get play => 'Jouer';

  @override
  String get difficulty => 'Difficulté';

  @override
  String get author => 'Auteur';

  @override
  String get clues => 'indices';

  @override
  String get difficultyEasy => 'Facile';

  @override
  String get difficultyMedium => 'Moyen';

  @override
  String get difficultyNormal => 'Normal';

  @override
  String get difficultyHard => 'Difficile';

  @override
  String get difficultyExpert => 'Expert';
}
