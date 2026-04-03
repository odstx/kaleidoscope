import SwiftUI

@main
struct FrontendApp: App {
    @StateObject private var authState = AuthState()
    @StateObject private var localization = Localization()
    
    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(authState)
                .environmentObject(localization)
        }
    }
}
