// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "FrontendApp",
    platforms: [
        .iOS(.v17),
        .macOS(.v14)
    ],
    products: [
        .executable(name: "FrontendApp", targets: ["FrontendApp"])
    ],
    targets: [
        .executableTarget(
            name: "FrontendApp",
            path: "Sources"),
        .testTarget(
            name: "FrontendAppTests",
            dependencies: ["FrontendApp"],
            path: "Tests/FrontendAppTests")
    ]
)
