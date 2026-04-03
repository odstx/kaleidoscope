import Foundation
import Combine

class AuthState: ObservableObject {
    @Published var isAuthenticated: Bool = false
    @Published var token: String? = nil
    
    private let tokenKey = "auth_token"
    
    init() {
        loadToken()
    }
    
    func login(token: String) {
        self.token = token
        self.isAuthenticated = true
        saveToken(token)
    }
    
    func logout() {
        self.token = nil
        self.isAuthenticated = false
        removeToken()
    }
    
    private func saveToken(_ token: String) {
        UserDefaults.standard.set(token, forKey: tokenKey)
    }
    
    private func loadToken() {
        if let savedToken = UserDefaults.standard.string(forKey: tokenKey) {
            self.token = savedToken
            self.isAuthenticated = true
        }
    }
    
    private func removeToken() {
        UserDefaults.standard.removeObject(forKey: tokenKey)
    }
}
