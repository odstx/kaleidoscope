import SwiftUI

struct DashboardView: View {
    @EnvironmentObject var authState: AuthState
    @EnvironmentObject var localization: Localization
    @EnvironmentObject var router: Router
    
    private var t: DashboardTranslations {
        localization.t.dashboard
    }
    
    var body: some View {
        VStack(spacing: 0) {
            NavbarView()
            
            ScrollView {
                VStack(alignment: .leading, spacing: 24) {
                    Text(t.title)
                        .font(.title)
                        .fontWeight(.bold)
                    
                    VStack(spacing: 16) {
                        CardView(title: t.welcome, subtitle: t.welcomeDesc) {
                            Text(t.welcomeMessage)
                                .foregroundColor(.secondary)
                        }
                    }
                }
                .padding()
            }
        }
    }
}

struct CardView<Content: View>: View {
    let title: String
    let subtitle: String
    @ViewBuilder let content: () -> Content
    
    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            VStack(alignment: .leading, spacing: 4) {
                Text(title)
                    .font(.headline)
                    .fontWeight(.semibold)
                
                Text(subtitle)
                    .font(.caption)
                    .foregroundColor(.secondary)
            }
            
            content()
        }
        .padding()
        .background(Color.secondary.opacity(0.1))
        .cornerRadius(12)
    }
}
