/// Application configuration.
class AppConfig {
  /// Base URL for the API server.
  final String apiBaseUrl;

  /// Default language code.
  final String defaultLanguage;

  const AppConfig({
    required this.apiBaseUrl,
    this.defaultLanguage = 'fr',
  });

  /// Development configuration.
  static const development = AppConfig(
    apiBaseUrl: 'http://localhost:8080',
  );

  /// Production configuration.
  static const production = AppConfig(
    apiBaseUrl: 'https://api.lesmotsdatche.com',
  );
}
