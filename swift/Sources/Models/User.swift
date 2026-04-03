import Foundation

struct User: Codable {
    let id: Int?
    let email: String
    let username: String?
}

struct LoginRequest: Codable {
    let email: String
    let password: String
}

struct RegisterRequest: Codable {
    let username: String
    let email: String
    let password: String
}

struct LoginResponse: Codable {
    let token: String
}

struct SystemInfo: Codable {
    let version: String
    let buildId: String
    let buildTime: String
    let gitCommit: String
    let openapiPath: String
    
    enum CodingKeys: String, CodingKey {
        case version
        case buildId = "build_id"
        case buildTime = "build_time"
        case gitCommit = "git_commit"
        case openapiPath = "openapi_path"
    }
}
