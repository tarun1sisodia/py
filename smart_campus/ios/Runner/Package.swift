// swift-tools-version:5.7
import PackageDescription

let package = Package(
    name: "SmartCampus",
    platforms: [
        .iOS(.v13)
    ],
    products: [
        .library(
            name: "SmartCampus",
            targets: ["SmartCampus"]),
    ],
    dependencies: [
        .package(
            url: "https://github.com/firebase/firebase-ios-sdk.git",
            .upToNextMajor(from: "10.0.0")
        ),
    ],
    targets: [
        .target(
            name: "SmartCampus",
            dependencies: [
                .product(name: "FirebaseAnalytics", package: "firebase-ios-sdk"),
                .product(name: "FirebaseAuth", package: "firebase-ios-sdk"),
                .product(name: "FirebaseDynamicLinks", package: "firebase-ios-sdk"),
                .product(name: "FirebaseFirestore", package: "firebase-ios-sdk")
            ]
        )
    ]
) 