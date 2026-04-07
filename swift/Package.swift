// swift-tools-version: 6.1
import PackageDescription

let package = Package(
    name: "Kaleidoscope",
    platforms: [
        .iOS(.v17),
        .macOS(.v14)
    ],
    products: [
        .executable(name: "Kaleidoscope", targets: ["Kaleidoscope"])
    ],
    targets: [
        .executableTarget(
            name: "Kaleidoscope",
            path: "Sources"),
        .testTarget(
            name: "KaleidoscopeTests",
            dependencies: ["Kaleidoscope"],
            path: "Tests/KaleidoscopeTests")
    ]
)
