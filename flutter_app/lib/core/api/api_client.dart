import 'package:dio/dio.dart';

import 'api_exceptions.dart';

/// HTTP client wrapper for API calls.
class ApiClient {
  final Dio _dio;

  ApiClient({required String baseUrl})
      : _dio = Dio(
          BaseOptions(
            baseUrl: baseUrl,
            connectTimeout: const Duration(seconds: 10),
            receiveTimeout: const Duration(seconds: 10),
            headers: {
              'Accept': 'application/json',
              'Accept-Encoding': 'gzip',
            },
          ),
        ) {
    _dio.interceptors.add(
      LogInterceptor(
        requestBody: true,
        responseBody: true,
        error: true,
      ),
    );
  }

  /// Performs a GET request.
  Future<Response<dynamic>> get(
    String path, {
    Map<String, dynamic>? queryParams,
    String? etag,
  }) async {
    try {
      final options = Options(
        headers: etag != null ? {'If-None-Match': etag} : null,
      );
      return await _dio.get(
        path,
        queryParameters: queryParams,
        options: options,
      );
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  ApiException _handleError(DioException e) {
    switch (e.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.sendTimeout:
      case DioExceptionType.receiveTimeout:
        return const NetworkException('Connection timeout');
      case DioExceptionType.connectionError:
        return const NetworkException('No internet connection');
      case DioExceptionType.badResponse:
        final statusCode = e.response?.statusCode;
        final message = _extractErrorMessage(e.response);
        if (statusCode == 404) {
          return NotFoundException(message);
        }
        if (statusCode == 401) {
          return UnauthorizedException(message);
        }
        return ServerException(message, statusCode: statusCode);
      default:
        return NetworkException(e.message ?? 'Unknown error');
    }
  }

  String _extractErrorMessage(Response<dynamic>? response) {
    if (response?.data is Map) {
      final data = response!.data as Map;
      return data['message'] as String? ??
          data['error'] as String? ??
          'Server error';
    }
    return 'Server error';
  }
}
