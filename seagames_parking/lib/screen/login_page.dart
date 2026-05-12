import 'package:flutter/material.dart';
import '../utils/app_theme.dart'; // เรียกใช้ Theme ที่เราสร้างไว้

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  // Controller สำหรับรับค่าจากช่องกรอก
  final _usernameController = TextEditingController();
  final _passwordController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppTheme.backgroundColor, // ใช้สีพื้นหลังจาก Theme
      body: SafeArea(
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.symmetric(horizontal: 30),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                // ส่วนของโลโก้หรือไอคอน
                const Icon(
                  Icons.local_parking_rounded,
                  size: 100,
                  color: AppTheme.primaryColor,
                ),
                const SizedBox(height: 20),
                const Text("Parking System", style: AppTheme.headingStyle),
                const SizedBox(height: 40),

                // ช่องกรอก Username
                _buildTextField(
                  controller: _usernameController,
                  hint: "Username",
                  icon: Icons.person_outline,
                ),
                const SizedBox(height: 15),

                // ช่องกรอก Password
                _buildTextField(
                  controller: _passwordController,
                  hint: "Password",
                  icon: Icons.lock_outline,
                  isPassword: true,
                ),
                const SizedBox(height: 30),

                // ปุ่มเข้าสู่ระบบ
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton(
                    onPressed: () {
                      // ตอนนี้ทำแค่ UI ให้กดแล้วไปหน้า Home ก่อน
                      // โดยส่ง Mock Role เป็น 'owner' ไปทดสอบ
                      Navigator.pushReplacementNamed(
                        context,
                        '/home',
                        // arguments: 'owner',
                      );
                    },
                    style: AppTheme.primaryButtonStyle, // ใช้สไตล์ปุ่มจาก Theme
                    child: const Text(
                      "LOGIN",
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  // Widget ตัวช่วยสำหรับสร้าง TextField เพื่อไม่ให้โค้ดรก
  Widget _buildTextField({
    required TextEditingController controller,
    required String hint,
    required IconData icon,
    bool isPassword = false,
  }) {
    return TextField(
      controller: controller,
      obscureText: isPassword,
      decoration: InputDecoration(
        prefixIcon: Icon(icon, color: AppTheme.primaryColor),
        hintText: hint,
        filled: true,
        fillColor: Colors.white,
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: BorderSide.none,
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: const BorderSide(color: Colors.black12),
        ),
      ),
    );
  }
}
