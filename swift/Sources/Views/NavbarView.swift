import SwiftUI

struct NavbarView: View {
    @EnvironmentObject var authState: AuthState
    @EnvironmentObject var localization: Localization
    @EnvironmentObject var router: Router
    @State private var showMenu = false
    @State private var showLanguageMenu = false
    
    private var t: NavTranslations {
        localization.t.nav
    }
    
    private var lang: LanguageTranslations {
        localization.t.languageLabel
    }
    
    var body: some View {
        HStack {
            Text(Bundle.main.object(forInfoDictionaryKey: "CFBundleName") as? String ?? "App")
                .font(.headline)
                .onTapGesture {
                    router.navigate(to: .dashboard)
                }
            
            Spacer()
            
            Menu {
                Button(action: { router.navigate(to: .profile) }) {
                    Label(t.profile, systemImage: "person.circle")
                }
                
                Menu {
                    Button(action: { localization.changeLanguage(.zh) }) {
                        Label(lang.zh, systemImage: localization.currentLanguage == .zh ? "checkmark" : "")
                    }
                    Button(action: { localization.changeLanguage(.en) }) {
                        Label(lang.en, systemImage: localization.currentLanguage == .en ? "checkmark" : "")
                    }
                } label: {
                    Label("\(lang.title) (\(localization.currentLanguage == .zh ? "中文" : "EN"))", systemImage: "globe")
                }
                
                Divider()
                
                Button(role: .destructive, action: { authState.logout() }) {
                    Label(t.logout, systemImage: "rectangle.portrait.and.arrow.right")
                }
            } label: {
                Text(t.menu)
                    .padding(.horizontal, 16)
                    .padding(.vertical, 8)
                    .background(Color.secondary.opacity(0.1))
                    .cornerRadius(8)
            }
        }
        .padding(.horizontal)
        .padding(.vertical, 12)
        .background(Color.primary.opacity(0.05))
        .shadow(color: Color.black.opacity(0.1), radius: 1, x: 0, y: 1)
    }
}
