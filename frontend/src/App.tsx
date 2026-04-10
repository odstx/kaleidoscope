import { BrowserRouter as Router, Routes, Route, Outlet } from 'react-router-dom';
import RegisterPage from './pages/RegisterPage';
import LoginPage from './pages/LoginPage';
import OIDCCallbackPage from './pages/OIDCCallbackPage';
import DashboardPage from './pages/DashboardPage';
import UserProfilePage from './pages/UserProfilePage';
import ForgotPasswordPage from './pages/ForgotPasswordPage';
import ResetPasswordPage from './pages/ResetPasswordPage';
import MicroAppPage from './pages/MicroAppPage';
import { Footer } from './components/Footer';
import { Navbar } from './components/Navbar';
import { AgentChat } from './components/AgentChat';
import { AuthProvider } from './contexts/AuthContext';
import { ThemeProvider } from './contexts/ThemeContext';
import './i18n';

type Layout = 'floating' | 'fullscreen';

function AppContent() {
  const handleLayoutChange = (_layout: Layout, _isOpen: boolean) => {
  };

  return (
    <div className="flex flex-col min-h-screen">
      <Navbar />
      <div className="flex flex-1">
        <div className="flex-1">
          <Outlet />
        </div>
      </div>
      <AgentChat onLayoutChange={handleLayoutChange} />
      <Footer />
    </div>
  );
}

function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <Router>
          <Routes>
            <Route path="/" element={<AppContent />}>
              <Route path="register" element={<RegisterPage />} />
              <Route path="login" element={<LoginPage />} />
              <Route path="oidc/callback" element={<OIDCCallbackPage />} />
              <Route path="forgot-password" element={<ForgotPasswordPage />} />
              <Route path="reset-password" element={<ResetPasswordPage />} />
              <Route path="dashboard" element={<DashboardPage />} />
              <Route path="profile" element={<UserProfilePage />} />
              <Route path="app/:appname" element={<MicroAppPage />} />
              <Route path="/" element={<LoginPage />} />
            </Route>
          </Routes>
        </Router>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;