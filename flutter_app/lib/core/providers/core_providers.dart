import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../api/api_client.dart';
import '../config/app_config.dart';

part 'core_providers.g.dart';

/// Provides the application configuration.
@riverpod
AppConfig appConfig(AppConfigRef ref) {
  // TODO: Switch based on environment
  return AppConfig.development;
}

/// Provides the API client.
@riverpod
ApiClient apiClient(ApiClientRef ref) {
  final config = ref.watch(appConfigProvider);
  return ApiClient(baseUrl: config.apiBaseUrl);
}
