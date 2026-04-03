import SwiftUI

struct LoginView: View {
    @EnvironmentObject var authState: AuthState
    @EnvironmentObject var localization: Localization
    @EnvironmentObject var router: Router
    
    @State private var email = ""
    @State private var password = ""
    @State private var isLoading = false
    @State private var errorMessage: String?
    
    private var t: LoginTranslations {
        localization.t.login
    }
    
    var body: some View {
        VStack(spacing: 20) {
            Spacer()
            
            VStack(spacing: 16) {
                Text(t.title)
                    .font(.largeTitle)
                    .fontWeight(.bold)
                
                Text(t.description)
                    .font(.subheadline)
                    .foregroundColor(.secondary)
            }
            
            VStack(spacing: 16) {
                if let error = errorMessage {
                    Text(error)
                        .foregroundColor(.red)
                        .padding()
                        .background(Color.red.opacity(0.1))
                        .cornerRadius(8)
                }
                
                VStack(alignment: .leading, spacing: 8) {
                    Text(t.email)
                        .font(.caption)
                        .foregroundColor(.secondary)
                    
                    TextField(t.emailPlaceholder, text: $email)
                        .textFieldStyle(.roundedBorder)
                        .textContentType(.emailAddress)
                    
                    if !email.isEmpty && !isValidEmail(email) {
                        Text(t.emailInvalid)
                            .font(.caption)
                            .foregroundColor(.red)
                    }
                }
                
                VStack(alignment: .leading, spacing: 8) {
                    Text(t.password)
                        .font(.caption)
                        .foregroundColor(.secondary)
                    
                    SecureField(t.passwordPlaceholder, text: $password)
                        .textFieldStyle(.roundedBorder)
                        .textContentType(.password)
                }
                
                Button(action: handleLogin) {
                    HStack {
                        if isLoading {
                            ProgressView()
                                .progressViewStyle(CircularProgressViewStyle())
                        }
                        Text(isLoading ? t.submitting : t.submit)
                    }
                    .frame(maxWidth: .infinity)
                    .padding()
                    .background(Color.accentColor)
                    .foregroundColor(.white)
                    .cornerRadius(8)
                }
                .disabled(isLoading || email.isEmpty || password.isEmpty || !isValidEmail(email))
            }
            .padding()
            
            Button(t.goToRegister) {
                router.navigate(to: .register)
            }
            .foregroundColor(.accentColor)
            
            Spacer()
        }
        .padding()
    }
    
    private func handleLogin() {
        guard isValidEmail(email) else { return }
        
        isLoading = true
        errorMessage = nil
        
        Task {
            do {
                let token = try await APIService.shared.login(email: email, password: password)
                await MainActor.run {
                    authState.login(token: token)
                    router.replace(with: .dashboard)
                }
            } catch {
                await MainActor.run {
                    errorMessage = error.localizedDescription
                    isLoading = false
                }
            }
        }
    }
    
    private func isValidEmail(_ email: String) -> Bool {
        let emailRegex = "[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,64}"
        return NSPredicate(format: "SELF MATCHES %@", emailRegex).evaluate(with: email)
    }
}
