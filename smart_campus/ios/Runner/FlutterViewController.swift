import SwiftUI
import Flutter

struct FlutterViewController: UIViewControllerRepresentable {
    func makeUIViewController(context: Context) -> FlutterViewController {
        return FlutterViewController(engine: FlutterEngine(name: "main"), nibName: nil, bundle: nil)
    }
    
    func updateUIViewController(_ uiViewController: FlutterViewController, context: Context) {
        // Updates can be handled here if needed
    }
} 