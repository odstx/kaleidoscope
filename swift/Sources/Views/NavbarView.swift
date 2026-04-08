import SwiftUI
import AppKit

struct NavbarView: View {
    @EnvironmentObject var authState: AuthState
    @EnvironmentObject var localization: Localization
    @EnvironmentObject var router: Router
    @EnvironmentObject var trayMenuManager: TrayMenuManager
    
    private var t: NavTranslations {
        localization.t.nav
    }
    
    var body: some View {
        HStack {
            Text(Bundle.main.object(forInfoDictionaryKey: "CFBundleName") as? String ?? "App")
                .font(.headline)
                .onTapGesture {
                    router.navigate(to: .dashboard)
                }
            
            Spacer()
        }
        .padding(.horizontal)
        .padding(.vertical, 12)
        .background(Color.primary.opacity(0.05))
        .shadow(color: Color.black.opacity(0.1), radius: 1, x: 0, y: 1)
        .onAppear {
            setupTrayMenu()
        }
    }
    
    private func setupTrayMenu() {
        trayMenuManager.setActions(
            logout: { authState.logout() },
            language: { showLanguageAlert() }
        )
    }
    
    private func showLanguageAlert() {
        let alert = NSAlert()
        alert.messageText = localization.t.languageLabel.title
        alert.addButton(withTitle: "中文")
        alert.addButton(withTitle: "English")
        
        if alert.runModal() == .alertFirstButtonReturn {
            localization.changeLanguage(.zh)
        } else {
            localization.changeLanguage(.en)
        }
    }
}
