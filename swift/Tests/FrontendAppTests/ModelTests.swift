import XCTest
@testable import FrontendApp

final class ModelTests: XCTestCase {
    func testUserInitialization() {
        let user = User(
            id: 1,
            email: "test@example.com",
            username: "testuser"
        )
        
        XCTAssertEqual(user.id, 1)
        XCTAssertEqual(user.username, "testuser")
        XCTAssertEqual(user.email, "test@example.com")
    }
    
    func testUserWithNilFields() {
        let user = User(
            id: nil,
            email: "test@example.com",
            username: nil
        )
        
        XCTAssertNil(user.id)
        XCTAssertNil(user.username)
        XCTAssertEqual(user.email, "test@example.com")
    }
    
    func testLoginRequestEncoding() {
        let request = LoginRequest(email: "test@example.com", password: "password123")
        
        let encoder = JSONEncoder()
        let data = try! encoder.encode(request)
        let json = try! JSONSerialization.jsonObject(with: data) as! [String: String]
        
        XCTAssertEqual(json["email"], "test@example.com")
        XCTAssertEqual(json["password"], "password123")
    }
    
    func testRegisterRequestEncoding() {
        let request = RegisterRequest(
            username: "testuser",
            email: "test@example.com",
            password: "password123"
        )
        
        let encoder = JSONEncoder()
        let data = try! encoder.encode(request)
        let json = try! JSONSerialization.jsonObject(with: data) as! [String: String]
        
        XCTAssertEqual(json["username"], "testuser")
        XCTAssertEqual(json["email"], "test@example.com")
        XCTAssertEqual(json["password"], "password123")
    }
    
    func testLoginResponseDecoding() {
        let json = """
        {
            "token": "test-token-123"
        }
        """.data(using: .utf8)!
        
        let decoder = JSONDecoder()
        let response = try! decoder.decode(LoginResponse.self, from: json)
        
        XCTAssertEqual(response.token, "test-token-123")
    }
    
    func testSystemInfoDecoding() {
        let json = """
        {
            "version": "1.0.0",
            "build_id": "build-123",
            "build_time": "2024-01-01T00:00:00Z",
            "git_commit": "abc123",
            "openapi_path": "/openapi.json"
        }
        """.data(using: .utf8)!
        
        let decoder = JSONDecoder()
        let info = try! decoder.decode(SystemInfo.self, from: json)
        
        XCTAssertEqual(info.version, "1.0.0")
        XCTAssertEqual(info.buildId, "build-123")
        XCTAssertEqual(info.buildTime, "2024-01-01T00:00:00Z")
        XCTAssertEqual(info.gitCommit, "abc123")
        XCTAssertEqual(info.openapiPath, "/openapi.json")
    }
    
    func testUserCodable() {
        let user = User(id: 42, email: "codable@test.com", username: "codableuser")
        
        let encoder = JSONEncoder()
        let data = try! encoder.encode(user)
        
        let decoder = JSONDecoder()
        let decodedUser = try! decoder.decode(User.self, from: data)
        
        XCTAssertEqual(decodedUser.id, 42)
        XCTAssertEqual(decodedUser.email, "codable@test.com")
        XCTAssertEqual(decodedUser.username, "codableuser")
    }
}
