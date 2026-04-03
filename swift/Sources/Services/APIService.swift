import Foundation

enum APIError: Error, LocalizedError {
    case invalidURL
    case httpError(statusCode: Int, message: String)
    case decodingError
    case encodingError
    case noData
    case unauthorized
    
    var errorDescription: String? {
        switch self {
        case .invalidURL:
            return "Invalid URL"
        case .httpError(let statusCode, let message):
            return "HTTP Error \(statusCode): \(message)"
        case .decodingError:
            return "Failed to decode response"
        case .encodingError:
            return "Failed to encode request"
        case .noData:
            return "No data received"
        case .unauthorized:
            return "Unauthorized"
        }
    }
}

class APIService {
    static let shared = APIService()
    
    private let baseURL: String
    private let session: URLSession
    
    private init() {
        self.baseURL = ProcessInfo.processInfo.environment["API_BASE_URL"] ?? "http://localhost:8000"
        self.session = URLSession.shared
    }
    
    func login(email: String, password: String) async throws -> String {
        let request = LoginRequest(email: email, password: password)
        let response: LoginResponse = try await post("/api/v1/users/login", body: request)
        return response.token
    }
    
    func register(username: String, email: String, password: String) async throws {
        let request = RegisterRequest(username: username, email: email, password: password)
        let _: LoginResponse = try await post("/api/v1/users/register", body: request)
    }
    
    func getUserInfo(token: String) async throws -> User {
        return try await get("/api/v1/users/info", token: token)
    }
    
    func getSystemInfo() async throws -> SystemInfo {
        return try await get("/api/v1/system/info")
    }
    
    private func get<T: Codable>(_ endpoint: String, token: String? = nil) async throws -> T {
        let url = try makeURL(endpoint)
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        if let token = token {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }
        
        return try await performRequest(request)
    }
    
    private func post<T: Codable, R: Codable>(_ endpoint: String, body: T, token: String? = nil) async throws -> R {
        let url = try makeURL(endpoint)
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        if let token = token {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }
        
        let encoder = JSONEncoder()
        encoder.keyEncodingStrategy = .convertToSnakeCase
        request.httpBody = try encoder.encode(body)
        
        return try await performRequest(request)
    }
    
    private func makeURL(_ endpoint: String) throws -> URL {
        guard let url = URL(string: baseURL + endpoint) else {
            throw APIError.invalidURL
        }
        return url
    }
    
    private func performRequest<T: Codable>(_ request: URLRequest) async throws -> T {
        let (data, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.noData
        }
        
        guard (200...299).contains(httpResponse.statusCode) else {
            let message = String(data: data, encoding: .utf8) ?? "Unknown error"
            if httpResponse.statusCode == 401 {
                throw APIError.unauthorized
            }
            throw APIError.httpError(statusCode: httpResponse.statusCode, message: message)
        }
        
        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        
        do {
            return try decoder.decode(T.self, from: data)
        } catch {
            throw APIError.decodingError
        }
    }
}
