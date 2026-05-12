import 'package:flutter/material.dart';
import 'package:salomon_bottom_bar/salomon_bottom_bar.dart';
import 'package:seagames_parking/utils/app_theme.dart';
import '../widgets/menu_card.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  int _selectedIndex = 0;

  late final List<Widget> _pages = [
    _buildHomeTab(),
    _buildHistoryTab(),
    const Center(child: Text("Search Page")),
    const Center(child: Text("Profile Page")),
  ];
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppTheme.backgroundColor,
      appBar: AppBar(
        title: const Text("SeaGames Parking"),
        backgroundColor: AppTheme.primaryColor,
        actions: [
          IconButton(
            icon: const Icon(Icons.logout),
            onPressed: () {
              // กลับไปหน้า Login
              Navigator.pushReplacementNamed(context, '/');
            },
          ),
        ],
      ),
      body: _pages[_selectedIndex],

      bottomNavigationBar: SalomonBottomBar(
        currentIndex: _selectedIndex,
        selectedItemColor: AppTheme.primaryColor,
        unselectedItemColor: const Color(0xff757575),
        onTap: (index) {
          setState(() {
            _selectedIndex = index;
          });
        },
        items: _navBarItems,
      ),
    );
  }
}

Widget _buildHomeTab() {
  return Padding(
    padding: const EdgeInsets.all(20.0),
    child: Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text("Menu", style: AppTheme.headerStyle),
        const SizedBox(height: 20),
        Expanded(
          child: GridView.count(
            crossAxisCount: 2,
            mainAxisSpacing: 15,
            crossAxisSpacing: 15,
            children: [
              MenuCard(
                title: "ออกตั๋วรถเข้า",
                icon: Icons.assignment_turned_in_rounded,
                iconColor: Colors.green,
                onTap: () {
                  /* TODO */
                },
              ),
              MenuCard(
                title: "สแกนรถออก",
                icon: Icons.qr_code_scanner_rounded,
                iconColor: Colors.orange,
                onTap: () {
                  // Navigator.pushNamed(context, '/scan');
                },
              ),
              // MenuCard(
              //   title: "ค้นหาข้อมูลรถ",
              //   icon: Icons.search_rounded,
              //   iconColor: Colors.blue,
              //   onTap: () {
              //     /* TODO */
              //   },
              // ),
              // MenuCard(
              //   title: "รายงานรายได้",
              //   icon: Icons.bar_chart_rounded,
              //   iconColor: Colors.purple,
              //   onTap: () {
              //     /* TODO */
              //   },
              // ),
            ],
          ),
        ),
      ],
    ),
  );
}

Widget _buildHistoryTab() {
  // ข้อมูลจำลอง (ในอนาคตจะดึงมาจาก Backend)
  final List<Map<String, String>> mockHistory = [
    {"plate": "กข 123", "in": "08:30", "out": "12:00", "status": "ออกแล้ว"},
    {"plate": "มล 999", "in": "09:15", "out": "-", "status": "กำลังจอด"},
    {"plate": "ชย 55", "in": "10:00", "out": "10:45", "status": "ออกแล้ว"},
  ];

  return Column(
    crossAxisAlignment: CrossAxisAlignment.start,
    children: [
      const Padding(
        padding: EdgeInsets.all(20.0),
        child: Text("ประวัติการเข้า-ออก", style: AppTheme.headerStyle),
      ),
      Expanded(
        child: ListView.builder(
          itemCount: mockHistory.length,
          itemBuilder: (context, index) {
            final item = mockHistory[index];
            return Card(
              margin: const EdgeInsets.symmetric(horizontal: 20, vertical: 8),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(12),
              ),
              child: ListTile(
                leading: CircleAvatar(
                  backgroundColor: item['status'] == "กำลังจอด"
                      ? Colors.green
                      : Colors.grey,
                  child: const Icon(Icons.motorcycle, color: Colors.white),
                ),
                title: Text(
                  "ทะเบียน: ${item['plate']}",
                  style: const TextStyle(fontWeight: FontWeight.bold),
                ),
                subtitle: Text("เข้า: ${item['in']} | ออก: ${item['out']}"),
                trailing: Text(
                  item['status']!,
                  style: TextStyle(
                    color: item['status'] == "กำลังจอด"
                        ? Colors.green
                        : Colors.black54,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
            );
          },
        ),
      ),
    ],
  );
}

final _navBarItems = [
  SalomonBottomBarItem(
    icon: const Icon(Icons.home),
    title: const Text("Home"),
    selectedColor: Colors.purple,
  ),
  SalomonBottomBarItem(
    icon: const Icon(Icons.history),
    title: const Text("History"),
    selectedColor: Colors.pink,
  ),
  SalomonBottomBarItem(
    icon: const Icon(Icons.search),
    title: const Text("Search"),
    selectedColor: Colors.orange,
  ),
  SalomonBottomBarItem(
    icon: const Icon(Icons.person),
    title: const Text("Profile"),
    selectedColor: Colors.teal,
  ),
];
