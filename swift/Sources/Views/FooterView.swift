import SwiftUI

struct FooterView: View {
    @EnvironmentObject var localization: Localization
    @State private var systemInfo: SystemInfo?
    @State private var isLoading = true
    @State private var hasError = false
    
    private var t: FooterTranslations {
        localization.t.footer
    }
    
    private var isDebugMode: Bool {
        #if DEBUG
        return true
        #else
        return false
        #endif
    }
    
    var body: some View {
        if isDebugMode {
            HStack(spacing: 16) {
                if hasError {
                    Text(t.error)
                        .font(.caption)
                        .foregroundColor(.secondary)
                } else if let info = systemInfo {
                    Text("\(t.version) **\(info.version)**")
                        .font(.caption)
                        .foregroundColor(.secondary)
                    
                    Text("\(t.buildId) **\(info.buildId)**")
                        .font(.caption)
                        .foregroundColor(.secondary)
                    
                    if let url = URL(string: "http://localhost:8000\(info.openapiPath)") {
                        Link("OpenAPI", destination: url)
                            .font(.caption)
                    }
                } else {
                    Text(t.loading)
                        .font(.caption)
                        .foregroundColor(.secondary)
                }
            }
            .padding(.vertical, 8)
            .frame(maxWidth: .infinity)
            .background(Color.secondary.opacity(0.1))
            .onAppear {
                loadSystemInfo()
            }
        }
    }
    
    private func loadSystemInfo() {
        Task {
            do {
                let info = try await APIService.shared.getSystemInfo()
                await MainActor.run {
                    self.systemInfo = info
                    self.isLoading = false
                }
            } catch {
                await MainActor.run {
                    self.hasError = true
                    self.isLoading = false
                }
            }
        }
    }
}
