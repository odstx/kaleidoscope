import SwiftUI

struct ProfileView: View {
    @EnvironmentObject var authState: AuthState
    @EnvironmentObject var localization: Localization
    @EnvironmentObject var router: Router
    
    @State private var user: User?
    @State private var isLoading = true
    @State private var errorMessage: String?
    
    private var t: ProfileTranslations {
        localization.t.profile
    }
    
    var body: some View {
        VStack(spacing: 0) {
            NavbarView()
            
            Group {
                if isLoading {
                    VStack {
                        Spacer()
                        Text(t.loading)
                            .foregroundColor(.secondary)
                        Spacer()
                    }
                } else if let error = errorMessage {
                    VStack {
                        Spacer()
                        Text(error)
                            .foregroundColor(.red)
                        Button("Retry") {
                            loadUserInfo()
                        }
                        .padding()
                        Spacer()
                    }
                } else if let user = user {
                    ScrollView {
                        VStack(alignment: .leading, spacing: 24) {
                            Text(t.title)
                                .font(.title)
                                .fontWeight(.bold)
                            
                            CardView(title: t.accountDetails, subtitle: t.accountDetailsDesc) {
                                VStack(spacing: 16) {
                                    HStack {
                                        Text(t.email)
                                            .foregroundColor(.secondary)
                                        Spacer()
                                        Text(user.email)
                                            .fontWeight(.medium)
                                    }
                                    
                                    if let userId = user.id {
                                        HStack {
                                            Text(t.userId)
                                                .foregroundColor(.secondary)
                                            Spacer()
                                            Text("\(userId)")
                                                .fontWeight(.medium)
                                        }
                                    }
                                }
                            }
                        }
                        .padding()
                    }
                }
            }
        }
        .onAppear {
            loadUserInfo()
        }
    }
    
    private func loadUserInfo() {
        guard let token = authState.token else {
            router.replace(with: .login)
            return
        }
        
        isLoading = true
        errorMessage = nil
        
        Task {
            do {
                let userInfo = try await APIService.shared.getUserInfo(token: token)
                await MainActor.run {
                    self.user = userInfo
                    self.isLoading = false
                }
            } catch {
                await MainActor.run {
                    self.errorMessage = t.fetchError
                    self.isLoading = false
                }
            }
        }
    }
}
