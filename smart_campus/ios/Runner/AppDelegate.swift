import SwiftUI
import Flutter
import FirebaseCore

class AppDelegate: NSObject, FlutterAppDelegate, UIApplicationDelegate {
    func application(
        _ application: UIApplication,
        didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]? = nil
    ) -> Bool {
        // Initialize Firebase
        FirebaseApp.configure()
        
        // Initialize Flutter
        GeneratedPluginRegistrant.register(with: self)
        
        return true
    }
}

@main
struct SmartCampusApp: App {
    // Register app delegate for Firebase setup
    @UIApplicationDelegateAdaptor(AppDelegate.self) var delegate
    
    var body: some Scene {
        WindowGroup {
            FlutterViewController()
                .edgesIgnoringSafeArea(.all)
        }
    }
}
