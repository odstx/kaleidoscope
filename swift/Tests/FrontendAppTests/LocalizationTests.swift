import XCTest
@testable import FrontendApp

final class LocalizationTests: XCTestCase {
    func testEnglishTranslations() {
        let localization = Localization()
        localization.changeLanguage(.en)
        
        XCTAssertEqual(localization.t.login.title, "User Login")
        XCTAssertEqual(localization.t.register.title, "User Registration")
        XCTAssertEqual(localization.t.dashboard.title, "Dashboard")
        XCTAssertEqual(localization.t.login.email, "Email")
        XCTAssertEqual(localization.t.login.password, "Password")
    }
    
    func testChineseTranslations() {
        let localization = Localization()
        localization.changeLanguage(.zh)
        
        XCTAssertEqual(localization.t.login.title, "用户登录")
        XCTAssertEqual(localization.t.register.title, "用户注册")
        XCTAssertEqual(localization.t.dashboard.title, "控制面板")
        XCTAssertEqual(localization.t.login.email, "邮箱")
        XCTAssertEqual(localization.t.login.password, "密码")
    }
    
    func testLanguageSwitch() {
        let localization = Localization()
        
        localization.changeLanguage(.en)
        XCTAssertEqual(localization.t.login.title, "User Login")
        
        localization.changeLanguage(.zh)
        XCTAssertEqual(localization.t.login.title, "用户登录")
        
        localization.changeLanguage(.en)
        XCTAssertEqual(localization.t.login.title, "User Login")
    }
}
