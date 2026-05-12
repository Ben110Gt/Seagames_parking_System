import 'package:flutter/material.dart';
import 'package:seagames_parking/screen/login_page.dart';
import 'package:seagames_parking/screen/home_sreen.dart';

void main() async {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'My Flutter App',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.blue),
        useMaterial3: true,
      ),
      home: const LoginScreen(),
      routes: {'/home': (context) => const HomeScreen()},
    );
  }
}
