import SwiftUI

struct RegisterView: View {
    @EnvironmentObject var localization: Localization
    @EnvironmentObject var router: Router
    
    @State private var username = ""
    @State private var email = ""
    @State private var password = ""
    @State private var isLoading = false
    @State private var errorMessage: String?
    @State private var showSuccessDialog = false
    @State private var countdown = 5
    
    private var t: RegisterTranslations {
        localization.t.register
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
                    Text(t.username)
                        .font(.caption)
                        .foregroundColor(.secondary)
                    
                    TextField(t.usernamePlaceholder, text: $username)
                        .textFieldStyle(.roundedBorder)
                        .textContentType(.username)
                    
                    if !username.isEmpty && !isValidUsername(username) {
                        Text(t.usernameInvalid)
                            .font(.caption)
                            .foregroundColor(.red)
                    }
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
                        .textContentType(.newPassword)
                    
                    if !password.isEmpty && password.count < 6 {
                        Text(t.passwordMin)
                            .font(.caption)
                            .foregroundColor(.red)
                    }
                }
                
                Button(action: handleRegister) {
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
                .disabled(isLoading || !isFormValid)
            }
            .padding()
            
            Button(t.goToLogin) {
                router.navigateBack()
            }
            .foregroundColor(.accentColor)
            
            Spacer()
        }
        .padding()
        .alert(t.successTitle, isPresented: $showSuccessDialog) {
            Button("OK") {
                router.replace(with: .login)
            }
        } message: {
            Text(String(format: t.successDescription, countdown))
        }
        .onChange(of: showSuccessDialog) { _, newValue in
            if newValue {
                startCountdown()
            }
        }
    }
    
    private var isFormValid: Bool {
        isValidUsername(username) && isValidEmail(email) && password.count >= 6
    }
    
    private func handleRegister() {
        guard isFormValid else { return }
        
        isLoading = true
        errorMessage = nil
        
        Task {
            do {
                try await APIService.shared.register(username: username, email: email, password: password)
                await MainActor.run {
                    isLoading = false
                    showSuccessDialog = true
                }
            } catch {
                await MainActor.run {
                    errorMessage = error.localizedDescription
                    isLoading = false
                }
            }
        }
    }
    
    private func startCountdown() {
        countdown = 5
        Timer.scheduledTimer(withTimeInterval: 1.0, repeats: true) { timer in
            countdown -= 1
            if countdown <= 0 {
                timer.invalidate()
                showSuccessDialog = false
                router.replace(with: .login)
            }
        }
    }
    
    private func isValidEmail(_ email: String) -> Bool {
        let emailRegex = "[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,64}"
        return NSPredicate(format: "SELF MATCHES %@", emailRegex).evaluate(with: email)
    }
    
    private func isValidUsername(_ username: String) -> Bool {
        guard username.count >= 3 && username.count <= 20 else { return false }
        let usernameRegex = "^[a-zA-Z0-9_]+$"
        return NSPredicate(format: "SELF MATCHES %@", usernameRegex).evaluate(with: username)
    }
}
