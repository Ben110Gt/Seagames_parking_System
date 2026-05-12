import 'package:flutter/material.dart';

class AppTheme {
  // กำหนดสีหลักของระบบจอดรถ
  static const Color primaryColor = Colors.blueAccent;
  static const Color accentColor = Colors.orangeAccent;
  static const Color backgroundColor = Color(0xFFF5F5F5);

  // กำหนดสไตล์ตัวอักษรหัวข้อ
  static const TextStyle headingStyle = TextStyle(
    fontSize: 24,
    fontWeight: FontWeight.bold,
    color: Colors.black87,
  );

  // กำหนดสไตล์ปุ่มเบื้องต้น
  static final ButtonStyle primaryButtonStyle = ElevatedButton.styleFrom(
    backgroundColor: primaryColor,
    foregroundColor: Colors.white,
    padding: const EdgeInsets.symmetric(vertical: 16),
    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
  );

  // สไตล์ตัวอักษรหัวข้อ
  static const TextStyle headerStyle = TextStyle(
    fontSize: 22,
    fontWeight: FontWeight.bold,
    color: Colors.black87,
  );
}
