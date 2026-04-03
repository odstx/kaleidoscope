import SwiftUI

enum Route: Hashable {
    case login
    case register
    case dashboard
    case profile
}

struct ContentView: View {
    @EnvironmentObject var authState: AuthState
    @State private var path = NavigationPath()
    
    var body: some View {
        NavigationStack(path: $path) {
            Group {
                if authState.isAuthenticated {
                    DashboardView()
                } else {
                    LoginView()
                }
            }
            .navigationDestination(for: Route.self) { route in
                switch route {
                case .login:
                    LoginView()
                case .register:
                    RegisterView()
                case .dashboard:
                    DashboardView()
                case .profile:
                    ProfileView()
                }
            }
        }
        .environmentObject(Router(path: $path))
    }
}

class Router: ObservableObject {
    @Published var path: NavigationPath
    
    init(path: Binding<NavigationPath>) {
        self._path = Published(initialValue: path.wrappedValue)
    }
    
    func navigate(to route: Route) {
        path.append(route)
    }
    
    func navigateBack() {
        path.removeLast()
    }
    
    func navigateToRoot() {
        path = NavigationPath()
    }
    
    func replace(with route: Route) {
        path = NavigationPath()
        path.append(route)
    }
}
