import Foundation

enum Language: String {
    case en = "en"
    case zh = "zh"
}

class Localization: ObservableObject {
    @Published var currentLanguage: Language {
        didSet {
            UserDefaults.standard.set(currentLanguage.rawValue, forKey: "language")
        }
    }
    
    init() {
        if let savedLang = UserDefaults.standard.string(forKey: "language"),
           let lang = Language(rawValue: savedLang) {
            self.currentLanguage = lang
        } else {
            self.currentLanguage = .zh
        }
    }
    
    func changeLanguage(_ language: Language) {
        currentLanguage = language
    }
    
    var t: Translations {
        Translations(language: currentLanguage)
    }
}

struct Translations {
    let language: Language
    
    var common: CommonTranslations {
        switch language {
        case .en:
            return CommonTranslations(
                loading: "Loading...",
                submit: "Submit",
                cancel: "Cancel",
                confirm: "Confirm"
            )
        case .zh:
            return CommonTranslations(
                loading: "加载中...",
                submit: "提交",
                cancel: "取消",
                confirm: "确认"
            )
        }
    }
    
    var nav: NavTranslations {
        switch language {
        case .en:
            return NavTranslations(
                dashboard: "Dashboard",
                menu: "Menu",
                profile: "Profile",
                logout: "Logout"
            )
        case .zh:
            return NavTranslations(
                dashboard: "控制台",
                menu: "菜单",
                profile: "个人信息",
                logout: "退出登录"
            )
        }
    }
    
    var login: LoginTranslations {
        switch language {
        case .en:
            return LoginTranslations(
                title: "User Login",
                description: "Enter your account information",
                email: "Email",
                emailPlaceholder: "Enter email",
                emailInvalid: "Please enter a valid email address",
                password: "Password",
                passwordPlaceholder: "Enter password",
                passwordRequired: "Please enter password",
                submit: "Login",
                submitting: "Logging in...",
                goToRegister: "Go to Register",
                loginFailed: "Login failed",
                loginError: "An error occurred during login"
            )
        case .zh:
            return LoginTranslations(
                title: "用户登录",
                description: "输入您的账户信息",
                email: "邮箱",
                emailPlaceholder: "请输入邮箱",
                emailInvalid: "请输入有效的邮箱地址",
                password: "密码",
                passwordPlaceholder: "请输入密码",
                passwordRequired: "请输入密码",
                submit: "登录",
                submitting: "登录中...",
                goToRegister: "去注册",
                loginFailed: "登录失败",
                loginError: "登录过程中发生错误"
            )
        }
    }
    
    var register: RegisterTranslations {
        switch language {
        case .en:
            return RegisterTranslations(
                title: "User Registration",
                description: "Create new account",
                username: "Username",
                usernamePlaceholder: "Enter username",
                usernameMin: "Username must be at least 3 characters",
                usernameMax: "Username must be at most 20 characters",
                usernameInvalid: "Username can only contain letters, numbers and underscores",
                email: "Email",
                emailPlaceholder: "Enter email",
                emailInvalid: "Please enter a valid email address",
                password: "Password",
                passwordPlaceholder: "Enter password",
                passwordMin: "Password must be at least 6 characters",
                submit: "Register",
                submitting: "Registering...",
                goToLogin: "Go to Login",
                registerFailed: "Registration failed",
                registerError: "An error occurred during registration",
                successTitle: "Registration Successful",
                successDescription: "Your account has been created successfully! Redirecting to login page in %d seconds"
            )
        case .zh:
            return RegisterTranslations(
                title: "用户注册",
                description: "创建新账户",
                username: "用户名",
                usernamePlaceholder: "请输入用户名",
                usernameMin: "用户名至少需要3个字符",
                usernameMax: "用户名最多20个字符",
                usernameInvalid: "用户名只能包含字母、数字和下划线",
                email: "邮箱",
                emailPlaceholder: "请输入邮箱",
                emailInvalid: "请输入有效的邮箱地址",
                password: "密码",
                passwordPlaceholder: "请输入密码",
                passwordMin: "密码至少需要6个字符",
                submit: "注册",
                submitting: "注册中...",
                goToLogin: "去登录",
                registerFailed: "注册失败",
                registerError: "注册过程中发生错误",
                successTitle: "注册成功",
                successDescription: "您的账户已成功创建！%d秒后自动跳转到登录页面"
            )
        }
    }
    
    var dashboard: DashboardTranslations {
        switch language {
        case .en:
            return DashboardTranslations(
                title: "Dashboard",
                quickActions: "Quick Actions",
                quickActionsDesc: "Common features",
                viewProfile: "View Profile",
                changePassword: "Change Password",
                welcome: "Welcome",
                welcomeDesc: "This is your personal dashboard",
                welcomeMessage: "You have successfully logged in. Here you can manage your personal information, view system status, and perform other operations."
            )
        case .zh:
            return DashboardTranslations(
                title: "控制面板",
                quickActions: "快速操作",
                quickActionsDesc: "常用功能入口",
                viewProfile: "查看个人信息",
                changePassword: "修改密码",
                welcome: "欢迎使用",
                welcomeDesc: "这是您的个人控制面板",
                welcomeMessage: "您已成功登录系统。在这里您可以管理您的个人信息、查看系统状态以及进行其他操作。"
            )
        }
    }
    
    var profile: ProfileTranslations {
        switch language {
        case .en:
            return ProfileTranslations(
                title: "User Profile",
                description: "Manage your account information",
                loading: "Loading...",
                accountDetails: "Account Details",
                accountDetailsDesc: "Your personal information",
                email: "Email:",
                userId: "User ID:",
                fetchError: "Failed to fetch user information"
            )
        case .zh:
            return ProfileTranslations(
                title: "用户信息",
                description: "管理您的账户信息",
                loading: "加载中...",
                accountDetails: "账户详情",
                accountDetailsDesc: "您的个人信息",
                email: "邮箱:",
                userId: "用户ID:",
                fetchError: "获取用户信息失败"
            )
        }
    }
    
    var footer: FooterTranslations {
        switch language {
        case .en:
            return FooterTranslations(
                loading: "Loading...",
                error: "Unable to load version info",
                version: "Version:",
                buildId: "Build ID:"
            )
        case .zh:
            return FooterTranslations(
                loading: "加载中...",
                error: "无法加载版本信息",
                version: "版本:",
                buildId: "构建ID:"
            )
        }
    }
    
    var languageLabel: LanguageTranslations {
        switch language {
        case .en:
            return LanguageTranslations(
                title: "Language",
                switchLanguage: "Switch Language",
                en: "English",
                zh: "中文"
            )
        case .zh:
            return LanguageTranslations(
                title: "语言",
                switchLanguage: "切换语言",
                en: "English",
                zh: "中文"
            )
        }
    }
}

struct CommonTranslations {
    let loading: String
    let submit: String
    let cancel: String
    let confirm: String
}

struct NavTranslations {
    let dashboard: String
    let menu: String
    let profile: String
    let logout: String
}

struct LoginTranslations {
    let title: String
    let description: String
    let email: String
    let emailPlaceholder: String
    let emailInvalid: String
    let password: String
    let passwordPlaceholder: String
    let passwordRequired: String
    let submit: String
    let submitting: String
    let goToRegister: String
    let loginFailed: String
    let loginError: String
}

struct RegisterTranslations {
    let title: String
    let description: String
    let username: String
    let usernamePlaceholder: String
    let usernameMin: String
    let usernameMax: String
    let usernameInvalid: String
    let email: String
    let emailPlaceholder: String
    let emailInvalid: String
    let password: String
    let passwordPlaceholder: String
    let passwordMin: String
    let submit: String
    let submitting: String
    let goToLogin: String
    let registerFailed: String
    let registerError: String
    let successTitle: String
    let successDescription: String
}

struct DashboardTranslations {
    let title: String
    let quickActions: String
    let quickActionsDesc: String
    let viewProfile: String
    let changePassword: String
    let welcome: String
    let welcomeDesc: String
    let welcomeMessage: String
}

struct ProfileTranslations {
    let title: String
    let description: String
    let loading: String
    let accountDetails: String
    let accountDetailsDesc: String
    let email: String
    let userId: String
    let fetchError: String
}

struct FooterTranslations {
    let loading: String
    let error: String
    let version: String
    let buildId: String
}

struct LanguageTranslations {
    let title: String
    let switchLanguage: String
    let en: String
    let zh: String
}
