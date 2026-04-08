import SwiftUI
import AppKit

@main
struct FrontendApp: App {
    @NSApplicationDelegateAdaptor(AppDelegate.self) var appDelegate
    @StateObject private var authState = AuthState()
    @StateObject private var localization = Localization()
    @StateObject private var trayMenuManager = TrayMenuManager()
    
    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(authState)
                .environmentObject(localization)
                .environmentObject(trayMenuManager)
                .onAppear {
                    trayMenuManager.isTrayVisible = true
                }
        }
    }
}

class AppDelegate: NSObject, NSApplicationDelegate {
    func applicationDidFinishLaunching(_ notification: Notification) {
        NSApp.setActivationPolicy(.regular)
    }
}

@MainActor
class TrayMenuManager: ObservableObject {
    @Published var statusItem: NSStatusItem?
    private var menu: NSMenu?
    
    @Published var isTrayVisible: Bool = false {
        didSet {
            if isTrayVisible {
                setupTray()
            } else {
                removeTray()
            }
        }
    }
    
    func setupTray() {
        guard statusItem == nil else { return }
        
        statusItem = NSStatusBar.system.statusItem(withLength: NSStatusItem.variableLength)
        
        if let button = statusItem?.button {
            button.image = NSImage(systemSymbolName: "line.3.horizontal", accessibilityDescription: "Menu")
        }
        
        menu = NSMenu()
        statusItem?.menu = menu
        
        buildMenu()
        
        isTrayVisible = true
    }
    
    func removeTray() {
        if let item = statusItem {
            NSStatusBar.system.removeStatusItem(item)
        }
        statusItem = nil
        menu = nil
    }
    
    private func buildMenu() {
        guard let menu = menu else { return }
        menu.removeAllItems()
        
        let languageItem = NSMenuItem(title: "Language", action: #selector(languageMenuAction), keyEquivalent: "l")
        languageItem.target = self
        languageItem.image = NSImage(systemSymbolName: "globe", accessibilityDescription: nil)
        menu.addItem(languageItem)
        
        menu.addItem(NSMenuItem.separator())
        
        let logoutItem = NSMenuItem(title: "Logout", action: #selector(logoutMenuAction), keyEquivalent: "q")
        logoutItem.target = self
        logoutItem.image = NSImage(systemSymbolName: "rectangle.portrait.and.arrow.right", accessibilityDescription: nil)
        menu.addItem(logoutItem)
    }
    
    private var logoutActionHandler: (@MainActor () -> Void)?
    private var languageActionHandler: (@MainActor () -> Void)?
    
    func setActions(
        logout: @escaping @MainActor () -> Void,
        language: @escaping @MainActor () -> Void
    ) {
        self.logoutActionHandler = logout
        self.languageActionHandler = language
    }
    
    @objc private func languageMenuAction() {
        languageActionHandler?()
    }
    
    @objc private func logoutMenuAction() {
        logoutActionHandler?()
    }
}
