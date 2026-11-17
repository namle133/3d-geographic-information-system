(function () {
  const STORAGE_KEY = 'vinhomes-language';
  const translations = {
    'BẢNG ĐIỀU KHIỂN GIS / BIM 3D': '3D GIS / BIM CONTROL PANEL',
    '01 CẮT LỚP CHIỀU CAO & THI CÔNG': '01 HEIGHT SECTION & CONSTRUCTION',
    '02 QUẢN LÝ LỚP DỮ LIỆU': '02 DATA LAYER MANAGEMENT',
    '03 THÔNG TIN PHÂN KHU CHỌN': '03 SELECTED AREA INFORMATION',
    'CẮT LỚP CHIỀU CAO & THI CÔNG': 'HEIGHT SECTION & CONSTRUCTION',
    '04 ĐIỀU CHỈNH ĐỘ CAO Z-OFFSET': '04 ADJUST HEIGHT Z-OFFSET',
    '05 TRẠNG THÁI CĂN HỘ (GIS)': '05 APARTMENT STATUS (GIS)',
    'Quản lý Yêu cầu Bảo trì': 'Maintenance Request Management',
    'Quản lý Tòa nhà — Tầng — Căn hộ': 'Building — Floor — Apartment Management',
    'Quản lý Tài khoản': 'Account Management',
    'Quản lý Hóa đơn': 'Invoice Management',
    'Quản lý Hợp đồng': 'Contract Management',
    'Quản lý Cư dân': 'Resident Management',
    'Báo cáo Tài chính': 'Financial Reports',
    'Vinhomes West Point - Mô hình 3D (ArcGIS)': 'Vinhomes West Point - 3D Model (ArcGIS)',
    'UC 2.7.46 — Tổng quan vận hành theo thời gian thực': 'UC 2.7.46 — Real-time Operations Overview',
    'Click vào mô hình 3D để quét và hiển thị dữ liệu kỹ thuật công trình': 'Click the 3D model to inspect and display building technical data',
    'Vinhomes Property Management': 'Vinhomes Property Management',
    'Cổng quản trị': 'Admin Portal',
    'Mở GIS 3D': 'Open 3D GIS',
    'Mở/đóng menu': 'Open/close menu',
    'Quên mật khẩu?': 'Forgot password?',
    'Ghi nhớ mật khẩu': 'Remember me',
    'Đăng nhập ngay': 'Log in',
    'Chưa có tài khoản?': "Don't have an account?",
    'Đăng ký ngay': 'Register now',
    'Đăng nhập với Google': 'Sign in with Google',
    'Đăng nhập với Facebook': 'Sign in with Facebook',
    'Email đăng ký': 'Registered email',
    'Gửi mã đặt lại': 'Send reset code',
    'Mật khẩu mới': 'New password',
    'Xác nhận mật khẩu': 'Confirm password',
    'Đặt lại mật khẩu': 'Reset password',
    'Sao chép': 'Copy',
    'Mã reset': 'Reset code',
    'Dashboard KPI': 'KPI Dashboard',
    'Tổng quan vận hành': 'Operations overview',
    'Hợp đồng sắp hết hạn (30 ngày)': 'Contracts expiring soon (30 days)',
    'Hóa đơn quá hạn': 'Overdue invoices',
    'Yêu cầu bảo trì chờ xử lý': 'Maintenance requests awaiting action',
    'Yêu cầu mới': 'New requests',
    'Đang xử lý': 'In progress',
    'Hoàn thành': 'Completed',
    'Khẩn cấp': 'Urgent',
    'Tất cả trạng thái': 'All statuses',
    'Tất cả mức độ ưu tiên': 'All priority levels',
    'Tất cả': 'All',
    'Từ ngày': 'From date',
    'Đến ngày': 'To date',
    'Tháng này': 'This month',
    'Tháng trước': 'Last month',
    'Quý này': 'This quarter',
    'Năm nay': 'This year',
    'Lọc hóa đơn': 'Filter invoices',
    'Tạo báo cáo': 'Generate report',
    'Xuất CSV': 'Export CSV',
    'Tạo yêu cầu bảo trì': 'Create maintenance request',
    'Danh sách Yêu cầu Bảo trì': 'Maintenance Requests',
    'Loại vấn đề': 'Issue type',
    'Mức độ ưu tiên': 'Priority',
    'Phân công cho nhân viên (tùy chọn)': 'Assign to staff (optional)',
    'Ghi chú phân công': 'Assignment notes',
    'Chi tiết yêu cầu': 'Request details',
    'Trạng thái mới': 'New status',
    'Lý do thay đổi': 'Reason for change',
    'Cập nhật trạng thái căn hộ': 'Update apartment status',
    'Tổng số tầng': 'Total floors',
    'Năm xây dựng': 'Year built',
    'Tìm tòa nhà...': 'Search buildings...',
    'Tìm mã căn hộ...': 'Search apartment number...',
    'Tìm kiếm theo căn hộ hoặc loại vấn đề...': 'Search by apartment or issue type...',
    'Tên tòa nhà': 'Building name',
    'Thêm tòa nhà': 'Add building',
    'Thêm cư dân': 'Add resident',
    'Thêm tài khoản': 'Add account',
    'Tạo hóa đơn': 'Create invoice',
    'Lưu hóa đơn': 'Save invoice',
    'Ngày tạo / Kỳ TT': 'Issue date / Billing period',
    'Các hạng mục': 'Line items',
    'Phương thức thanh toán': 'Payment method',
    'Tổng căn hộ': 'Total apartments',
    'Toàn khu Vinhomes West Point': 'Vinhomes West Point complex',
    'Đang thuê': 'Rented',
    'Đang trống': 'Vacant',
    'Sẵn sàng cho thuê': 'Available for rent',
    'Đang bảo trì': 'Under maintenance',
    'Cần theo dõi': 'Needs attention',
    'Doanh thu tháng': 'Monthly revenue',
    'Từ HĐ active · Đã thu:': 'From active contracts · Collected:',
    'Sắp hết HĐ': 'Expiring soon',
    'Mở GIS 3D →': 'Open 3D GIS →',
    'Tòa nhà': 'Buildings',
    'Căn hộ': 'Apartments',
    'Tầng': 'Floor',
    'Tổng quan': 'Overview',
    'Hợp đồng': 'Contracts',
    'Hóa đơn': 'Invoices',
    'Bảo trì': 'Maintenance',
    'Báo cáo': 'Reports',
    'Tài khoản': 'Accounts',
    'Thông báo': 'Notifications',
    'Quản trị viên': 'Administrator',
    'Ban quản lý': 'Management',
    'Nhân viên': 'Staff',
    'Đăng xuất': 'Log out',
    'Menu': 'Menu',
    'Đăng nhập': 'Log in',
    'Đăng ký': 'Register',
    'Thanh toán': 'Payment',
    'Kết quả thanh toán': 'Payment result',
    'Tài khoản của tôi': 'My account',
    'Hồ sơ cá nhân': 'Profile',
    'Tên đăng nhập': 'Username',
    'Họ tên': 'Full name',
    'Họ tên *': 'Full name *',
    'Điện thoại': 'Phone',
    'Điện thoại *': 'Phone *',
    'Số CCCD/CMND': 'ID number',
    'Vai trò': 'Role',
    'Trạng thái': 'Status',
    'Địa chỉ': 'Address',
    'Mô tả': 'Description',
    'Ghi chú': 'Notes',
    'Lưu': 'Save',
    'Hủy': 'Cancel',
    'Đóng': 'Close',
    'Xóa': 'Delete',
    'Sửa': 'Edit',
    'Xem': 'View',
    'Chi tiết': 'Details',
    'Cập nhật': 'Update',
    'Thêm': 'Add',
    'Tạo mới': 'Create',
    'Tìm kiếm': 'Search',
    'Tìm': 'Search',
    'Đổi TT': 'Change status',
    'Thấp': 'Low',
    'Trung bình': 'Medium',
    'Cao': 'High',
    'Mới': 'New',
    'Đã phân công': 'Assigned',
    'Hủy bỏ': 'Cancelled',
    'Trống': 'Vacant',
    'Đã thuê': 'Rented',
    'Chưa có căn hộ': 'No apartments yet',
    'Không có hóa đơn quá hạn': 'No overdue invoices',
    'Không có yêu cầu chờ xử lý': 'No requests awaiting action',
    'Không có hợp đồng sắp hết hạn': 'No contracts expiring soon',
    'Không có dữ liệu': 'No data available',
    'Chưa có dữ liệu': 'No data available yet',
    'Đã cập nhật trạng thái — GIS 3D sẽ đổi màu': 'Status updated — the 3D GIS will update its colors',
    'Cuộn xuống danh sách quá hạn': 'Scroll to the overdue list',
    'Danh sách căn hộ': 'Apartment list',
    'Danh sách tòa nhà': 'Building list',
    'Danh sách cư dân': 'Resident list',
    'Danh sách hợp đồng': 'Contract list',
    'Danh sách hóa đơn': 'Invoice list',
    'Danh sách thông báo': 'Notification list',
    'Chi tiết căn hộ': 'Apartment details',
    'Chi tiết hóa đơn': 'Invoice details',
    'Chi tiết hợp đồng': 'Contract details',
    'Thông tin căn hộ': 'Apartment information',
    'Thông tin cư dân': 'Resident information',
    'Thông tin tài khoản': 'Account information',
    'Thông tin thanh toán': 'Payment information',
    'Thông tin chung': 'General information',
    'Thông tin liên hệ': 'Contact information',
    'Mã căn hộ': 'Apartment number',
    'Mã tòa': 'Building code',
    'Số tầng nổi': 'Above-ground floors',
    'Chiều cao (m)': 'Height (m)',
    'Phân bổ tầng': 'Floor distribution',
    'Tổng số căn': 'Total apartments',
    'Hết HĐ': 'Contract ends',
    'Ngày hết hạn': 'Expiration date',
    'Ngày bắt đầu': 'Start date',
    'Ngày kết thúc': 'End date',
    'Ngày tạo': 'Created date',
    'Ngày thanh toán': 'Payment date',
    'Kỳ thanh toán': 'Billing period',
    'Số tiền': 'Amount',
    'Tổng tiền': 'Total amount',
    'Đơn giá': 'Unit price',
    'Số lượng': 'Quantity',
    'Thành tiền': 'Line total',
    'Đã thanh toán': 'Paid',
    'Chưa thanh toán': 'Unpaid',
    'Quá hạn': 'Overdue',
    'Đang hoạt động': 'Active',
    'Không hoạt động': 'Inactive',
    'Đã hủy': 'Cancelled',
    'Đã hoàn thành': 'Completed',
    'Người thực hiện': 'Assigned technician',
    'Người tạo': 'Created by',
    'Người nhận': 'Recipient',
    'Nội dung': 'Content',
    'Loại thông báo': 'Notification type',
    'Đánh dấu đã đọc': 'Mark as read',
    'Đánh dấu tất cả đã đọc': 'Mark all as read',
    'Chưa đọc': 'Unread',
    'Đã đọc': 'Read',
    'Không tìm thấy': 'Not found',
    'Không tìm thấy dữ liệu': 'No data found',
    'Vui lòng chọn': 'Please select',
    'Vui lòng nhập': 'Please enter',
    'Bắt buộc': 'Required',
    'Tùy chọn': 'Optional',
    'Đang cập nhật...': 'Updating...',
    'Đang lưu...': 'Saving...',
    'Đang xử lý...': 'Processing...',
    'Tải lại dữ liệu': 'Reload data',
    'Làm mới dữ liệu': 'Refresh data',
    'Chọn tòa nhà': 'Select a building',
    'Chọn căn hộ': 'Select an apartment',
    'Chọn nhân viên': 'Select a staff member',
    'Chọn trạng thái': 'Select a status',
    'Chọn phương thức thanh toán': 'Select a payment method',
    'Cấp nước / Thoát nước': 'Plumbing / Drainage',
    'Điện / Điện tử': 'Electrical / Electronics',
    'Kết cấu / Vách tường': 'Structure / Walls',
    'Thiết bị / Nội thất': 'Appliances / Furniture',
    'Vệ sinh / Dọn dẹp': 'Cleaning',
    'Email không hợp lệ': 'Invalid email',
    'Đăng nhập thành công!': 'Login successful!',
    'Đăng ký thành công!': 'Registration successful!',
    'Đăng nhập thất bại': 'Login failed',
    'Đăng ký thất bại': 'Registration failed',
    'Đặt lại mật khẩu thành công. Vui lòng đăng nhập bằng mật khẩu mới.': 'Password reset successfully. Please sign in with your new password.',
    'Bạn không có quyền truy cập trang này': 'You are not authorized to access this page',
    'Ở lại': 'Stay',
    'Chọn bản đồ nền': 'Choose a basemap',
    'Chiều cao:': 'Height:',
    'Tòa nhà:': 'Building:',
    'Tầng:': 'Floor:',
    'Căn hộ:': 'Apartment:',
    'Trạng thái:': 'Status:',
    'Hết HĐ:': 'Contract ends:',
    'Nền sân dự án (Xám đậm)': 'Project ground (dark gray)',
    'Đường Phạm Hùng': 'Pham Hung Road',
    'Sân chơi trẻ em trung tâm': "Central children's playground",
    'Đường dạo chính tầng mái': 'Main rooftop walkway',
    'Bể bơi bốn mùa': 'All-season swimming pool',
    'Mái kính bể bơi': 'Pool glass canopy',
    'Khu thiền yoga': 'Yoga and meditation area',
    'Khu ngắm cảnh thành phố': 'City viewing area',
    'Ghế nghỉ': 'Rest bench',
    'Người dùng': 'Users',
    'Cư dân': 'Residents',
    'Tạo yêu cầu': 'Create request',
    'Tạo hợp đồng': 'Create contract',
    'Tạo tài khoản': 'Create account',
    'Tạo cư dân': 'Create resident',
    'Cập nhật thông tin': 'Update information',
    'Cập nhật trạng thái': 'Update status',
    'Xác nhận': 'Confirm',
    'Xác nhận thanh toán': 'Confirm payment',
    'Thanh toán thành công': 'Payment successful',
    'Thanh toán thất bại': 'Payment failed',
    'Đang chờ xử lý': 'Pending',
    'Đã xử lý': 'Processed',
    'Yêu cầu bảo trì': 'Maintenance request',
    'Hợp đồng thuê': 'Lease agreement',
    'Bảng điều khiển': 'Dashboard',
    'Tổng số': 'Total',
    'Tổng doanh thu': 'Total revenue',
    'Tổng hóa đơn': 'Total invoices',
    'Tổng cư dân': 'Total residents',
    'Tổng hợp đồng': 'Total contracts',
    'Tháng': 'Month',
    'Năm': 'Year',
    'Ngày': 'Date',
    'Mã': 'Code',
    'Loại': 'Type',
    'Tên': 'Name',
    'Vị trí': 'Location',
    'Tiền mặt': 'Cash',
    'Chuyển khoản': 'Bank transfer',
    'Thẻ': 'Card',
    'Trở về': 'Back',
    'Quay lại': 'Back',
    'Tiếp tục': 'Continue',
    'Tải xuống': 'Download',
    'Xem thêm': 'See more',
    'Chưa có thông báo': 'No notifications yet',
    'Chưa có hợp đồng': 'No contracts yet',
    'Chưa có hóa đơn': 'No invoices yet',
    'Chưa có cư dân': 'No residents yet',
    'Chưa có yêu cầu bảo trì': 'No maintenance requests yet',
    'Lỗi khi tải dữ liệu': 'Error loading data',
    'Vui lòng thử lại': 'Please try again',
    'HOẶC': 'OR',
    'Đăng nhập -': 'Log in -',
    'Đăng ký -': 'Register -',
    'Quên mật khẩu -': 'Forgot password -',
    'Kết quả thanh toán -': 'Payment result -',
    'Có': 'Yes',
    'Không': 'No',
    'Đang tải dữ liệu...': 'Loading data...',
    'Đang tải...': 'Loading...',
    'Thử lại': 'Try again',
    'Xem tất cả': 'View all',
    'Xem & thu tiền →': 'View & collect →',
    'Phạm Hùng, Mễ Trì, Nam Từ Liêm, Hà Nội': 'Pham Hung, Me Tri, Nam Tu Liem, Hanoi',
    'Mặt cắt chiều cao:': 'Height section:',
    '0m (Móng)': '0m (Foundation)',
    '160m (Ăng-ten)': '160m (Antenna)',
    '► Quét thi công tự động': '► Automatic construction scan',
    'Reset về mặc định': 'Reset to default',
    'Nền & Đường': 'Ground & Roads',
    'Khối đế': 'Podium',
    'Khung tháp': 'Tower Frame',
    'Vỏ kính': 'Glass Facade',
    'Đỉnh tháp': 'Tower Top',
    'Tiện ích mái': 'Rooftop Amenities',
    'Cảnh quan': 'Landscaping',
    'Ánh sáng': 'Lighting',
    'Địa hình & Nền móng': 'Terrain & Foundation',
    'Nền & Địa hình:': 'Ground & Terrain:',
    'Chế độ Ban đêm (Night)': 'Night mode',
    'Dữ liệu backend:': 'Backend data:',
    'Hiển thị lớp trạng thái': 'Show status layer',
    'Đường xá & Lối đi': 'Roads & Walkways',
    'Cảnh quan Sân vườn': 'Garden Landscaping',
    'Khối đế Thương mại': 'Commercial Podium',
    'Kết cấu cốt đế': 'Podium Structure',
    'Mặt vách kính đế': 'Podium Glass Facade',
    'Sàn mái tầng 4': 'Level 4 Roof Deck',
    '3 Tòa tháp căn hộ': '3 Residential Towers',
    'Tháp West 1 (Hỗn hợp)': 'West Tower 1 (Mixed Use)',
    'Tháp West 2 (Căn hộ)': 'West Tower 2 (Residential)',
    'Tháp West 3 (Căn hộ)': 'West Tower 3 (Residential)',
    'Hệ ban công căn hộ': 'Apartment Balcony System',
    'Thanh nẹp nhôm vỏ kính': 'Glass Facade Aluminum Trim',
    'Tiện ích Tầng mái (18m)': 'Rooftop Amenities (18m)',
    'Hệ thống Cây xanh': 'Greenery System',
    'Trạng thái căn hộ': 'Apartment Status',
    'Trạng thái GIS': 'GIS Status',
    'Click để quay về hướng Bắc': 'Click to face north',
    'Nhập email': 'Enter email',
    'Mật khẩu': 'Password',
    'Email': 'Email',
    'English': 'English',
    'Tiếng Việt': 'Vietnamese',
    'Chờ duyệt': 'Pending approval',
    'Gia hạn': 'Renew',
    'Kết thúc': 'Terminate',
    '-- Chọn cư dân --': '-- Select resident --',
    '-- Chọn tài khoản --': '-- Select account --',
    '-- Chọn phương thức --': '-- Select payment method --',
    '-- Chọn loại vấn đề --': '-- Select issue type --',
    '-- Chưa phân công --': '-- Unassigned --',
    'Tiền thuê/tháng (VND) *': 'Monthly rent (VND) *',
    'Tiền cọc (VND) *': 'Security deposit (VND) *',
    'Điều khoản và điều kiện': 'Terms and conditions',
    'Xử lý hợp đồng': 'Manage contract',
    'Nhập các điều khoản và điều kiện của hợp đồng...': 'Enter the terms and conditions of the contract...',
    'Duyệt hợp đồng': 'Approve contract',
    'Gia hạn hợp đồng': 'Extend contract',
    'Quên mật khẩu': 'Forgot password',
    'Bước 1 / 2': 'Step 1 / 2',
    'Bước 2 / 2': 'Step 2 / 2',
    'Nhập mã reset và mật khẩu mới (tối thiểu 6 ký tự).': 'Enter the reset code and new password (at least 6 characters).',
    'Tối thiểu 6 ký tự': 'At least 6 characters',
    'Đã sao chép mã reset': 'Reset code copied',
    'Đã gửi mã reset. Kiểm tra bước 2 bên dưới.': 'Reset code sent. Check step 2 below.',
    'Không thể sao chép — hãy chọn và copy thủ công': 'Unable to copy — select and copy it manually',
    'Nháp': 'Draft',
    'Chờ thanh toán': 'Awaiting payment',
    'Đã phát hành': 'Issued',
    'Hạn thanh toán': 'Due date',
    'Ghi nhận thanh toán': 'Record payment',
    'Số hóa đơn': 'Invoice number',
    'Séc': 'Check',
    'Tiền thuê căn hộ': 'Apartment rent',
    'Vui lòng điền đầy đủ thông tin bắt buộc': 'Please complete all required fields',
    'Tải ảnh (không bắt buộc)': 'Upload photos (optional)',
    '📷 Nhấp để chọn ảnh hoặc kéo thả ảnh vào đây': '📷 Click to select photos or drag and drop them here',
    '📷 Nhấp để chọn ảnh': '📷 Click to select photos',
    'Phân công yêu cầu': 'Assign request',
    'Phân công': 'Assign',
    'Chi phí (VND)': 'Cost (VND)',
    'Ảnh hoàn thành công việc': 'Work completion photos',
    'UC 2.7.45 — Quản lý thông báo hệ thống': 'UC 2.7.45 — System Notification Management',
    '+ Soạn thông báo': '+ Compose notification',
    'Hệ thống (Admin)': 'System (Admin)',
    'Soạn thông báo': 'Compose notification',
    'Gửi tới': 'Send to',
    'Toàn bộ người dùng': 'All users',
    'Một người cụ thể': 'A specific user',
    'Tiêu đề': 'Title',
    'Đã nhận thanh toán': 'Payment received',
    'Gửi thông báo': 'Send notification',
    'Đang xác minh thanh toán...': 'Verifying payment...',
    'Đổi mật khẩu': 'Change password',
    'Xác Thực': 'Verification',
    'Mã Xác Thực': 'Verification code',
    'Mã được cấp khi ký hợp đồng': 'Code provided when signing the contract',
    'Bạn sẽ nhận được mã này từ phòng bán hàng': 'You will receive this code from the sales office',
    'Bạn sẽ nhận được mã này từ phòng bán hàng': 'You will receive this code from the sales office',
    'Tài Khoản': 'Account',
    'Mật Khẩu': 'Password',
    'Xác Nhận Mật Khẩu': 'Confirm password',
    'Thông Tin Cá Nhân': 'Personal information',
    'CCCD/CMND': 'National ID',
    'Ngày Sinh': 'Date of birth',
    'Giới Tính': 'Gender',
    'Nam': 'Male',
    '-- Chọn giới tính --': '-- Select gender --',
    'Nữ': 'Female',
    'Khác': 'Other',
    'Số Điện Thoại': 'Phone number',
    'Địa Chỉ Thường Trú': 'Permanent address',
    'Đăng Ký': 'Register',
    'Đã có tài khoản?': 'Already have an account?',
    'UC 2.7.39 — Tổng hợp doanh thu & công nợ': 'UC 2.7.39 — Revenue and Debt Summary',
    'Chờ TT': 'Awaiting payment',
    'Chọn khoảng thời gian và nhấn': 'Select a date range and click',
    'Doanh thu theo kỳ thanh toán': 'Revenue by billing period',
    'Danh sách công nợ (hiện tại)': 'Current outstanding balances',
    'Gửi bảo trì': 'Submit maintenance request',
    'UC 2.7.7 — Admin quản lý người dùng hệ thống': 'UC 2.7.7 — System User Management',
    'Hoạt động': 'Active',
    'Vô hiệu': 'Disabled',
    'Khóa': 'Lock',
    'Vỉa hè': 'Sidewalk',
    'Lối đi': 'Walkway',
    '⏸ Tạm dừng quét': '⏸ Pause scan',
    'Trung tâm thương mại Vincom, siêu thị Vinmart/Vinmart+, shophouse chân đế, sảnh Lounge, phòng gym': 'Vincom mall, Vinmart/Vinmart+ supermarket, podium shophouses, lounge lobby, and gym',
    'Bể bơi mái kính, vườn nổi, khu vui chơi, không gian cảnh quan': 'Glass-roof pool, elevated garden, play area, and landscaped space',
    'Kết nối 3 tầng hầm để xe của dự án': 'Connected to the project’s three basement parking levels',
    'Phục vụ nhu cầu mua sắm, sinh hoạt, giải trí cho cư dân W1, W2, W3': 'Serves the shopping, daily living, and leisure needs of W1, W2, and W3 residents',
    'Khối ngang lớn liên kết cả 3 tòa, mái làm cảnh quan.': 'A large horizontal structure connecting all three towers, with a landscaped roof.',
    'Tòa West 3 - Vinhomes West Point': 'West Tower 3 - Vinhomes West Point',
    'Tòa West 2 - Vinhomes West Point': 'West Tower 2 - Vinhomes West Point',
    'Tòa West 1 - Vinhomes West Point': 'West Tower 1 - Vinhomes West Point',
    '14 căn hộ/mặt sàn': '14 apartments per floor',
    '7 thang thường + 1 thang hàng': '7 passenger elevators + 1 service elevator',
    'Liền tường': 'Built-in',
    'Sàn Footprint West 3': 'West 3 footprint slab',
    'Cột chịu lực Đông Bắc': 'Northeast support column',
    'Cột chịu lực Tây Bắc': 'Northwest support column',
    'Cột chịu lực Đông Nam': 'Southeast support column',
    'Cột chịu lực Tây Nam': 'Southwest support column',
    'Lõi thang máy trung tâm khối đế': 'Podium central elevator core',
    'Mặt vách kính khối đế thương mại': 'Commercial podium glass facade',
    'Sàn mái khối đế': 'Podium roof slab',
    'Lan can mái khối đế': 'Podium roof railing',
    'Sàn tầng 2': 'Level 2 slab',
    'Sàn tầng 3': 'Level 3 slab',
    'Sàn tầng 4': 'Level 4 slab',
    'Sàn tầng 5': 'Level 5 slab',
    'Siêu thị WinMart': 'WinMart supermarket',
    'Cửa hàng Thời trang Zara': 'Zara fashion store',
    'Cửa hàng Apple Store': 'Apple Store',
    'Nhà hàng Haidilao Hotpot': 'Haidilao Hotpot restaurant',
    "Nhà hàng Pizza 4P's": "Pizza 4P's restaurant",
    'Khu vui chơi trẻ em TiniWorld': 'TiniWorld children’s play area',
    'Rạp chiếu phim CGV Cinema': 'CGV Cinema',
    'Trung tâm California Fitness': 'California Fitness center',
    'Gian hàng tại tầng ': 'Shop on level ',
    ' của Trung tâm thương mại': ' of the shopping mall',
    'Khách mua sắm': 'Shopper',
    'Nhóm bạn trẻ': 'Group of young friends',
    'Khách xem đồ Zara': 'Zara customer',
    'Khách trải nghiệm iPhone': 'Customer trying an iPhone',
    'Khách dạo chơi': 'Visitor',
    'Khách chờ bàn Haidilao': 'Guest waiting for a table at Haidilao',
    'Khách ăn Pizza': 'Guest having pizza',
    'Gia đình đi chơi': 'Family outing',
    'Nhóm trẻ em': 'Group of children',
    'Khách mua vé xem phim': 'Moviegoer',
    'Khách tập Gym': 'Gym member',
    'Khách tham quan': 'Visitor',
    'đang tham quan mua sắm tại Vinhomes West Point.': ' is shopping at Vinhomes West Point.',
    'Đường Đỗ Đức Dục': 'Do Duc Duc Road',
    'Sảnh đón trả khách Phạm Hùng': 'Pham Hung passenger drop-off and pick-up lobby',
    'Bể phun nước tiểu cảnh': 'Decorative fountain ',
    'Thảm cỏ tròn tiểu cảnh': 'Decorative circular lawn ',
    'Cột đèn thép': 'Steel lamppost ',
    'Cần đèn': 'Lamp arm ',
    'Bóng đèn': 'Light bulb ',
    'Thảm cỏ trung tâm': 'Central lawn',
    'Sân chơi trẻ em': 'Children’s playground',
    'căn hộ': 'apartment',
    'Diện tích (m²)': 'Area (m²)',
    'Bạn có chắc muốn đăng xuất khỏi hệ thống? Phiên làm việc hiện tại sẽ kết thúc.': 'Are you sure you want to log out? Your current session will end.',
    'Bạn có chắc chắn muốn tiếp tục?': 'Are you sure you want to continue?',
    'Mật khẩu phải ít nhất 6 ký tự': 'Password must be at least 6 characters',
    'Đăng nhập thành công! Chuyển hướng...': 'Login successful! Redirecting...',
    'Lỗi kết nối server': 'Server connection error',
    'Tính năng này sắp có mặt': 'This feature is coming soon',
    'Đặt trước': 'Reserved',
    'Đã thêm tòa nhà': 'Building added',
    'Lỗi tải dữ liệu:': 'Error loading data:',
    'Lỗi tải hóa đơn:': 'Error loading invoices:',
    'Lỗi tải cư dân:': 'Error loading residents:',
    'Lỗi tải căn hộ:': 'Error loading apartments:',
    'Lỗi tải yêu cầu bảo trì:': 'Error loading maintenance requests:',
    'Lỗi tạo yêu cầu bảo trì': 'Error creating maintenance request',
    'Yêu cầu đã được phân công thành công': 'Request assigned successfully',
    'Lỗi phân công yêu cầu': 'Error assigning request',
    'Tiêu đề thông báo...': 'Notification title...',
    'HĐ sắp hết hạn': 'Contract expiring soon',
    'HĐ đã duyệt': 'Contract approved',
    'HĐ bị từ chối': 'Contract rejected',
    'Phân công bảo trì': 'Maintenance assignment',
    'Đã đánh dấu đọc': 'Marked as read',
    'Đã xóa': 'Deleted',
    'Đã đánh dấu tất cả đã đọc': 'Marked all as read',
    'Đã lưu thông tin (cần API PUT /profile trên backend)': 'Information saved (requires the backend PUT /profile API)',
    'Đã cập nhật mật khẩu (cần API trên backend)': 'Password updated (requires a backend API)',
    'Nhập lại mật khẩu': 'Re-enter password',
    'Nhập địa chỉ đầy đủ (số nhà, đường, quận, thành phố)...': 'Enter the full address (house number, street, district, city)...',
    'CCCD/CMND phải từ 9-12 chữ số': 'National ID must contain 9–12 digits',
    'Số điện thoại không hợp lệ (tối thiểu 10 chữ số)': 'Invalid phone number (at least 10 digits)',
    'Phải từ 18 tuổi trở lên': 'Must be at least 18 years old',
    'Ngày kết thúc phải sau ngày bắt đầu': 'End date must be after the start date',
    'Đang tạo...': 'Generating...',
    'Đang tổng hợp dữ liệu...': 'Compiling data...',
    'Báo cáo đã được tạo': 'Report generated',
    'Lỗi tải KPI:': 'Error loading KPIs:',
    'Không tải được KPI:': 'Unable to load KPIs:',
    'Không tải hóa đơn quá hạn:': 'Unable to load overdue invoices:',
    'Đã cập nhật tài khoản': 'Account updated',
    'Đã tạo tài khoản': 'Account created',
    'Vô hiệu hóa tài khoản này?': 'Deactivate this account?',
    'Đã vô hiệu hóa': 'Deactivated',
    'Lỗi:': 'Error:',
    'Ngôn ngữ': 'Language',
    'Quản lý Tòa nhà - Vinhomes': 'Building Management - Vinhomes',
    '+ Thêm hợp đồng': '+ Add contract',
    'Tìm kiếm theo số hợp đồng hoặc tên cư dân...': 'Search by contract number or resident name...',
    'Thêm hợp đồng': 'Add contract',
    'Kỳ thanh toán (tháng) *': 'Billing period (months) *',
    'Khôi phục tài khoản Vinhomes Property Management': 'Recover your Vinhomes Property Management account',
    'Nhập email đã đăng ký. Hệ thống sẽ tạo mã đặt lại mật khẩu (môi trường demo hiển thị mã ngay trên màn hình).': 'Enter your registered email. The system will create a password reset code (in the demo, the code appears on this screen).',
    'Mã 32 ký tự': '32-character code',
    'vd: admin@vinhomes.com': 'e.g. admin@vinhomes.com',
    'Quay lại đăng nhập': 'Back to log in',
    'Tìm kiếm số hóa đơn hoặc tên cư dân...': 'Search by invoice number or resident name...',
    '+ Thêm hạng mục': '+ Add line item',
    'Ghi chú hay mô tả về hóa đơn...': 'Notes or invoice description...',
    'Sửa hóa đơn': 'Edit invoice',
    'Số tiền phải thanh toán': 'Amount due',
    'Thanh toán trực tuyến': 'Online payment',
    'Số tiền thanh toán': 'Payment amount',
    'Mã tham chiếu': 'Reference number',
    'Mô tả chi tiết': 'Detailed description',
    'Mô tả vấn đề, triệu chứng, vị trí cụ thể...': 'Describe the issue, symptoms, and exact location...',
    'Thêm ghi chú cho nhân viên...': 'Add notes for staff...',
    'Hoàn thành yêu cầu': 'Complete request',
    'Mô tả giải pháp': 'Describe the resolution',
    'Mô tả những gì đã được làm để giải quyết vấn đề...': 'Describe the work done to resolve the issue...',
    'Ghi chú hoàn thành': 'Completion notes',
    'Ghi chú bổ sung...': 'Additional notes...',
    'Tên nhân viên...': 'Staff member name...',
    'Tất cả cư dân': 'All residents',
    'Tất cả nhân viên': 'All staff',
    'Nội dung chi tiết...': 'Detailed message...',
    'Hợp đồng sắp hết hạn': 'Contracts expiring soon',
    'Cập nhật bảo trì': 'Maintenance update',
    'UC 2.7.5 / 2.7.6 — Xem và chỉnh sửa thông tin cá nhân': 'UC 2.7.5 / 2.7.6 — View and edit personal information',
    'Mật khẩu hiện tại': 'Current password',
    'Xác nhận mật khẩu mới': 'Confirm new password',
    'Cập nhật mật khẩu': 'Update password',
    'Đăng ký tài khoản cư dân': 'Register a resident account',
    'Mã được cấp khi ký hợp đồng': 'Code provided when signing the contract',
    'Tên Đăng Nhập': 'Username',
    'Họ Tên Đầy Đủ': 'Full name',
    'Tất cả trong kỳ': 'All in period',
    'Hóa đơn trong kỳ báo cáo': 'Invoices in the reporting period',
    'UC 2.7.47 — Tổng quan cá nhân': 'UC 2.7.47 — Personal overview',
    'Hóa đơn chờ thanh toán': 'Invoices awaiting payment',
    'Không có hóa đơn chờ thanh toán': 'No invoices awaiting payment',
    'Yêu cầu bảo trì đang xử lý': 'Maintenance requests in progress',
    'Không có yêu cầu bảo trì': 'No maintenance requests',
    'Tìm kiếm cư dân...': 'Search residents...',
    'Tất cả vai trò': 'All roles',
    'Hợp đồng của tôi': 'My contracts',
    'Tìm kiếm theo số hợp đồng...': 'Search by contract number...',
    'Vui lòng chọn cư dân': 'Please select a resident',
    'Vui lòng chọn ngày bắt đầu': 'Please select a start date',
    'Vui lòng chọn ngày kết thúc': 'Please select an end date',
    'Vui lòng nhập tiền thuê hợp lệ': 'Please enter a valid monthly rent',
    'Kỳ thanh toán phải ≥ 1 tháng': 'Billing period must be at least 1 month',
    'Vui lòng nhập tiền cọc hợp lệ': 'Please enter a valid deposit',
    'Không tìm thấy hợp đồng': 'Contract not found',
    'Kết thúc hợp đồng': 'Terminate contract',
    'Cập nhật hợp đồng thành công': 'Contract updated successfully',
    'Tạo hợp đồng thành công': 'Contract created successfully',
    'Duyệt hợp đồng thành công': 'Contract approved successfully',
    'Vui lòng chọn ngày kết thúc mới': 'Please select a new end date',
    'Gia hạn hợp đồng thành công': 'Contract extended successfully',
    'Kết thúc hợp đồng thành công': 'Contract terminated successfully',
    'Mật khẩu xác nhận không khớp': 'The passwords do not match',
    'Mật khẩu phải có ít nhất 6 ký tự': 'Password must be at least 6 characters',
    'Đặt lại mật khẩu thành công! Đang chuyển đến trang đăng nhập...': 'Password reset successfully! Redirecting to the login page...',
    'Đặt lại mật khẩu thất bại': 'Password reset failed',
    '— Chưa có cư dân có HĐ active —': '— No residents with an active contract —',
    '— Không có căn hộ —': '— No apartments available —',
    'Hóa đơn đã được cập nhật': 'Invoice updated',
    'Hóa đơn đã được tạo thành công': 'Invoice created successfully',
    'Thanh toán đã được ghi nhận': 'Payment recorded',
    '— Chưa có căn hộ thuê (cần HĐ active) —': '— No rented apartments (requires an active contract) —',
    'Yêu cầu bảo trì đã được tạo thành công': 'Maintenance request created successfully',
    'Vui lòng chọn nhân viên': 'Please select a staff member',
    'Yêu cầu bảo trì đã được hoàn thành': 'Maintenance request completed',
    'Hóa đơn mới': 'New invoice',
    'Bảo trì thiết bị': 'Equipment maintenance',
    'Vui lòng chọn người nhận': 'Please select a recipient',
    'Ngày sinh không hợp lệ': 'Invalid date of birth',
    'Ngày sinh không thể là ngày tương lai': 'Date of birth cannot be in the future',
    'Họ tên phải có ít nhất 2 ký tự': 'Full name must be at least 2 characters',
    'Họ tên không quá 100 ký tự': 'Full name must be no more than 100 characters',
    'Tên đăng nhập phải có ít nhất 3 ký tự': 'Username must be at least 3 characters',
    'Vui lòng chọn giới tính': 'Please select a gender',
    'Địa chỉ thường trú phải có ít nhất 5 ký tự': 'Permanent address must be at least 5 characters',
    'Đăng ký thành công! Chuyển hướng đến trang đăng nhập...': 'Registration successful! Redirecting to the login page...',
    'Cư dân,Căn hộ,Số HĐ nợ,Quá hạn,Tổng nợ': 'Residents,Apartment,Outstanding invoices,Overdue,Total outstanding',
    'Vui lòng chọn khoảng thời gian': 'Please select a date range',
    'Xóa cư dân thành công': 'Resident deleted',
    'Cập nhật cư dân thành công': 'Resident updated',
    'Thêm cư dân thành công': 'Resident added',
    'Sửa tài khoản': 'Edit account',
    'Khối đế thương mại - tiện ích Vinhomes West Point': 'Commercial podium and amenities - Vinhomes West Point',
    'Khối đế thương mại, dịch vụ và tiện ích nội khu': 'Commercial podium, services, and on-site amenities',
    'Tầng 1-5: Shophouse, sảnh Lounge, shop thương mại, Officetel': 'Levels 1-5: Shophouses, lounge lobby, retail stores, and officetel units',
    'Tầng 1: Sảnh Lounge & shop thương mại; Tầng 2-35: Căn hộ chung cư': 'Level 1: Lounge lobby and retail shops; Levels 2-35: Apartments',
    'Tầng 1: Sảnh Lounge & shop thương mại; Tầng 2-4: Trung tâm thương mại & Officetel; Tầng 5-37: Căn hộ chung cư; Tầng 38-39: Căn hộ Duplex': 'Level 1: Lounge lobby and retail shops; Levels 2-4: Shopping mall and officetel units; Levels 5-37: Apartments; Levels 38-39: Duplex apartments',
    'Tầng 5-18: 12 căn hộ/mặt sàn; Tầng 19-39: 13 căn hộ/mặt sàn': 'Levels 5-18: 12 apartments per floor; Levels 19-39: 13 apartments per floor',
    'Tầng 1: Sảnh Lounge & shop thương mại; Tầng 2-4: Trung tâm thương mại & Officetel; Tầng 5-37: Officetel; Tầng 38-39: Căn hộ Duplex': 'Level 1: Lounge lobby and retail shops; Levels 2-4: Shopping mall and officetel units; Levels 5-37: Officetel units; Levels 38-39: Duplex apartments',
    'Tầng 5-37: 11 căn hộ/mặt sàn; Tầng 38-39: 13 căn hộ/mặt sàn': 'Levels 5-37: 11 apartments per floor; Levels 38-39: 13 apartments per floor',
    'Sàn Footprint Khối đế chính West 1 & 2': 'Main West 1 & 2 podium footprint slab',
    'Sảnh Lounge Khối đế': 'Podium lounge lobby',
    'Thang cuốn Tầng 1 -> Tầng 2': 'Escalator Level 1 -> Level 2',
    'Thang cuốn Tầng 2 -> Tầng 3': 'Escalator Level 2 -> Level 3',
    'Thang cuốn Tầng 3 -> Tầng 4': 'Escalator Level 3 -> Level 4',
    'Mật khẩu phải có ít nhất 8 ký tự': 'Password must be at least 8 characters',
    'Vui lòng điền đầy đủ thông tin': 'Please complete all required information',
    'Vui lòng thêm ít nhất một hạng mục với giá > 0': 'Add at least one line item with a price greater than 0',
    'Lỗi lưu hóa đơn': 'Error saving invoice',
    'Còn thiếu ': 'Remaining ',
    '✓ Đủ': '✓ Paid in full',
    'Vui lòng điền đầy đủ thông tin thanh toán': 'Please complete all payment information',
    'Lỗi ghi nhận thanh toán': 'Error recording payment',
    'Đang chuyển đến VNPay...': 'Redirecting to VNPay...',
    'Không nhận được URL thanh toán': 'No payment URL was received',
    'Không thể khởi tạo thanh toán VNPay': 'Unable to initiate VNPay payment',
    'Bạn có chắc chắn muốn hủy hóa đơn này?': 'Are you sure you want to cancel this invoice?',
    'Hạn thanh toán phải sau ngày tạo hóa đơn': 'The due date must be after the invoice date',
    'Ghi chú về thanh toán...': 'Payment notes...',
    'Lịch sử thanh toán': 'Payment history',
    'Hạn TT': 'Due date',
    'Tổng phát hành': 'Total issued',
    'Đã thu': 'Collected',
    'Còn nợ': 'Outstanding',
    'VD: Phòng khách, nhà bếp, ban công...': 'e.g. Living room, kitchen, balcony...',
    'Vui lòng mô tả giải pháp': 'Please describe the resolution',
    'Lỗi hoàn thành yêu cầu': 'Error completing request',
    'Đã chuyển sang trạng thái đang xử lý': 'Request moved to in-progress',
    'Bạn có chắc chắn muốn hủy yêu cầu này?': 'Are you sure you want to cancel this request?',
    'Yêu cầu đã được hủy': 'Request cancelled',
    'Đã gửi ': 'Sent ',
    ' thông báo': ' notification(s)',
    'Tối thiểu 8 ký tự': 'At least 8 characters',
    'Nguyễn Văn A': 'e.g. Alex Morgan',
    'Kỳ báo cáo:': 'Reporting period:',
    'BÁO CÁO TÀI CHÍNH - Vinhomes West Point': 'FINANCIAL REPORT - Vinhomes West Point',
    'Tổng phát hành,': 'Total issued,',
    'Đã thu,': 'Collected,',
    'Còn nợ kỳ,': 'Period outstanding,',
    'Công nợ hiện tại,': 'Current outstanding,',
    'CÔNG NỢ THEO CƯ DÂN': 'OUTSTANDING BY RESIDENT',
    'THEO KỲ': 'BY BILLING PERIOD',
    'Kỳ,Số HĐ,Phát hành,Đã thu,Còn nợ': 'Period,Invoices,Issued,Collected,Outstanding',
    'Hiển thị tối đa 100 hóa đơn mới nhất trong kỳ.': 'Showing up to 100 of the most recent invoices in this period.',
    'Chỉnh sửa cư dân': 'Edit resident',
    'Bạn có chắc chắn muốn xóa cư dân này?': 'Are you sure you want to delete this resident?',
    'Xóa thông báo này?': 'Delete this notification?',
    '(để trống nếu không đổi)': '(leave blank to keep unchanged)',
    'Vách kỹ thuật tháp ': 'Tower technical shaft ',
    'Hóa đơn đã được hủy': 'Invoice cancelled',
    'Không thể gửi yêu cầu': 'Unable to send the request'
  };
  const entries = Object.entries(translations).sort((a, b) => b[0].length - a[0].length);
  const reverseTranslations = new Map();
  entries.forEach(([source, target]) => {
    if (!reverseTranslations.has(target)) reverseTranslations.set(target, source);
  });
  Object.entries({
    'Log in': 'Đăng nhập',
    'Completed': 'Đã hoàn thành',
    'Total apartments': 'Tổng số căn',
    'Rented': 'Đang thuê',
    'Vacant': 'Đang trống',
    'Search': 'Tìm kiếm',
    'Cancelled': 'Đã hủy',
    'Back': 'Quay lại',
    'Awaiting payment': 'Chờ thanh toán',
    'Active': 'Đang hoạt động',
    'Pending': 'Đang chờ xử lý',
    'Paid': 'Đã thanh toán',
    'Unpaid': 'Chưa thanh toán',
    'Overdue': 'Quá hạn',
    'Approved': 'Đã duyệt',
    'Rejected': 'Bị từ chối',
    'Unknown': 'Không xác định',
    'Server connection error': 'Lỗi kết nối máy chủ',
    'Invalid credentials': 'Thông tin đăng nhập không chính xác',
    'Unauthorized': 'Không có quyền truy cập',
    'Not found': 'Không tìm thấy',
    'Network Error': 'Lỗi kết nối mạng',
    'Failed to fetch': 'Không thể kết nối đến máy chủ',
    'Property Management': 'Quản lý bất động sản',
    'Email hoặc mật khẩu không đúng': 'Email or password is incorrect',
    'Email hoặc username là bắt buộc': 'Email or username is required',
    'Tài khoản đã bị khóa hoặc vô hiệu hóa': 'Account is locked or disabled',
    'Failed to generate token': 'Không thể tạo token đăng nhập',
    'Failed to hash password': 'Không thể mã hóa mật khẩu',
    'Apartment not found': 'Không tìm thấy căn hộ',
    'Resident not found': 'Không tìm thấy cư dân',
    'Contract not found': 'Không tìm thấy hợp đồng',
    'Failed to create contract': 'Không thể tạo hợp đồng',
    'Failed to extend contract': 'Không thể gia hạn hợp đồng',
    'Failed to update apartment status': 'Không thể cập nhật trạng thái căn hộ',
    'Maintenance request not found': 'Không tìm thấy yêu cầu bảo trì',
    'Failed to create maintenance request': 'Không thể tạo yêu cầu bảo trì',
    'Maintenance request assigned successfully': 'Yêu cầu bảo trì đã được phân công',
    'Maintenance request completed successfully': 'Yêu cầu bảo trì đã hoàn thành',
    'Failed to update status': 'Không thể cập nhật trạng thái',
    'Failed to fetch apartments': 'Không thể tải danh sách căn hộ',
    'Failed to fetch residents': 'Không thể tải danh sách cư dân',
    'Failed to fetch contracts': 'Không thể tải danh sách hợp đồng',
    'Failed to fetch invoices': 'Không thể tải danh sách hóa đơn',
    'Failed to fetch notifications': 'Không thể tải thông báo',
    'Failed to create notification': 'Không thể tạo thông báo',
    'Failed to mark notification as read': 'Không thể đánh dấu thông báo đã đọc',
    'Failed to mark notifications as read': 'Không thể đánh dấu các thông báo đã đọc',
    'Failed to delete notification': 'Không thể xóa thông báo',
    'Failed to create invoice': 'Không thể tạo hóa đơn',
    'Failed to record payment': 'Không thể ghi nhận thanh toán',
    'Invalid amount': 'Số tiền không hợp lệ',
    'Invalid signature': 'Chữ ký không hợp lệ',
    'Failed to fetch equipment': 'Không thể tải thiết bị',
    'Failed to register equipment': 'Không thể đăng ký thiết bị',
    'Failed to update maintenance record': 'Không thể cập nhật hồ sơ bảo trì',
    'Forbidden': 'Bạn không có quyền thực hiện thao tác này',
    'User not found': 'Không tìm thấy người dùng',
    'Failed to create building': 'Không thể tạo tòa nhà',
    'Failed to fetch buildings': 'Không thể tải danh sách tòa nhà',
    'Building not found': 'Không tìm thấy tòa nhà',
    'Failed to update building': 'Không thể cập nhật tòa nhà',
    'Failed to delete building': 'Không thể xóa tòa nhà',
    'Failed to create floor': 'Không thể tạo tầng',
    'Failed to create apartment': 'Không thể tạo căn hộ',
    'Failed to update apartment': 'Không thể cập nhật căn hộ',
    'Building updated successfully': 'Cập nhật tòa nhà thành công',
    'Building deleted successfully': 'Xóa tòa nhà thành công',
    'Apartment status updated successfully': 'Cập nhật trạng thái căn hộ thành công',
    'Failed to approve contract': 'Không thể duyệt hợp đồng',
    'Failed to terminate contract': 'Không thể kết thúc hợp đồng',
    'Contract extended successfully': 'Gia hạn hợp đồng thành công',
    'Contract terminated successfully': 'Kết thúc hợp đồng thành công',
    'Invoice not found': 'Không tìm thấy hóa đơn',
    'Failed to update invoice': 'Không thể cập nhật hóa đơn',
    'Failed to fetch overdue invoices': 'Không thể tải hóa đơn quá hạn',
    'Failed to fetch maintenance requests': 'Không thể tải yêu cầu bảo trì',
    'Failed to assign maintenance request': 'Không thể phân công yêu cầu bảo trì',
    'Failed to complete request': 'Không thể hoàn thành yêu cầu',
    'Equipment not found': 'Không tìm thấy thiết bị',
    'Status updated successfully': 'Cập nhật trạng thái thành công',
    'Maintenance record updated successfully': 'Cập nhật hồ sơ bảo trì thành công',
    'Failed to fetch audit logs': 'Không thể tải nhật ký kiểm tra',
    'Notification marked as read': 'Đã đánh dấu thông báo đã đọc',
    'All notifications marked as read': 'Đã đánh dấu tất cả thông báo đã đọc',
    'Notification deleted': 'Đã xóa thông báo',
    'Failed to create user': 'Không thể tạo người dùng',
    'Failed to create resident profile': 'Không thể tạo hồ sơ cư dân',
    'Failed to update profile': 'Không thể cập nhật hồ sơ',
    'Failed to add dependent': 'Không thể thêm người phụ thuộc',
    'Failed to remove dependent': 'Không thể xóa người phụ thuộc',
    'Failed to register vehicle': 'Không thể đăng ký phương tiện',
    'Failed to fetch vehicles': 'Không thể tải danh sách phương tiện',
    'Failed to unregister vehicle': 'Không thể hủy đăng ký phương tiện',
    'Resident created successfully': 'Tạo cư dân thành công',
    'Profile updated successfully': 'Cập nhật hồ sơ thành công',
    'Dependent removed successfully': 'Xóa người phụ thuộc thành công',
    'Vehicle unregistered successfully': 'Hủy đăng ký phương tiện thành công',
    'Failed to fetch users': 'Không thể tải danh sách người dùng',
    'Failed to deactivate user': 'Không thể vô hiệu hóa người dùng',
    'Low': 'Thấp',
    'Medium': 'Trung bình',
    'High': 'Cao',
    'Urgent': 'Khẩn cấp',
    'low': 'Thấp',
    'medium': 'Trung bình',
    'high': 'Cao',
    'urgent': 'Khẩn cấp',
    'plumbing': 'Cấp nước / Thoát nước',
    'electrical': 'Điện / Điện tử',
    'structural': 'Kết cấu / Vách tường',
    'appliances': 'Thiết bị / Nội thất',
    'cleaning': 'Vệ sinh / Dọn dẹp',
    'Contract Expiring Soon': 'Hợp đồng sắp hết hạn',
    'New Invoice Generated': 'Hóa đơn mới',
    'Overdue Payment Alert': 'Cảnh báo thanh toán quá hạn',
    'New Maintenance Request': 'Yêu cầu bảo trì mới',
    'Maintenance Request Assigned': 'Yêu cầu bảo trì được phân công',
    'Equipment Maintenance Due': 'Đến hạn bảo trì thiết bị',
    'Payment Successful': 'Thanh toán thành công'
  }).forEach(([source, target]) => reverseTranslations.set(source, target));
  const originals = new WeakMap();
  let language = localStorage.getItem(STORAGE_KEY) === 'en' ? 'en' : 'vi';
  let applying = false;
  const notificationPatterns = [
    [/^Your rental contract for apartment (.+) will expire on (.+)$/, (_, apartment, date) =>
      `Hợp đồng thuê căn hộ ${apartment} sẽ hết hạn vào ${date}`],
    [/^Invoice (.+) for amount (.+) VND is now due on (.+)$/, (_, invoice, amount, date) =>
      `Hóa đơn ${invoice} trị giá ${amount} VND đến hạn thanh toán vào ${date}`],
    [/^Your invoice (.+) is now overdue\. Outstanding amount: (.+) VND$/, (_, invoice, amount) =>
      `Hóa đơn ${invoice} đã quá hạn. Số tiền còn nợ: ${amount} VND`],
    [/^New ([^ ]+) maintenance request from apartment ([^:]+): (.*)$/, (_, priority, apartment, description) =>
      `Yêu cầu bảo trì mới (${reverseTranslations.get(priority) || priority}) từ căn hộ ${apartment}: ${description}`],
    [/^Maintenance request for apartment (.+) has been assigned to you\. Issue: (.+)$/, (_, apartment, issue) =>
      `Yêu cầu bảo trì của căn hộ ${apartment} đã được phân công cho bạn. Vấn đề: ${reverseTranslations.get(issue) || issue}`],
    [/^Equipment '(.+)' \((.+)\) is due for maintenance$/, (_, name, type) =>
      `Thiết bị '${name}' (${type}) đã đến hạn bảo trì`]
  ];

  function translateText(value) {
    if (language === 'en') {
      return entries.reduce((text, [source, target]) => text.split(source).join(target), value);
    }
    for (const [pattern, format] of notificationPatterns) {
      const match = value.match(pattern);
      if (match) return format(...match);
    }
    return reverseTranslations.get(value) || value;
  }

  function translateNode(node) {
    const prior = originals.get(node);
    const source = prior && node.nodeValue === prior.translated ? prior.source : node.nodeValue;
    const translated = translateText(source);
    originals.set(node, { source, translated });
    if (node.nodeValue !== translated) node.nodeValue = translated;
  }

  function translateElement(element) {
    if (element.nodeType === Node.TEXT_NODE) {
      translateNode(element);
      return;
    }
    if (element.nodeType !== Node.ELEMENT_NODE) return;
    if (element.tagName === 'SCRIPT' || element.tagName === 'STYLE' || element.tagName === 'NOSCRIPT') return;

    const translatableAttributes = ['placeholder', 'title', 'aria-label', 'alt'];
    if (element.matches('input[type="button"], input[type="submit"], input[type="reset"]')) {
      translatableAttributes.push('value');
    }
    for (const attr of translatableAttributes) {
      if (!element.hasAttribute(attr)) continue;
      let values = originals.get(element);
      if (!values) {
        values = {};
        originals.set(element, values);
      }
      const current = element.getAttribute(attr);
      const previous = values[attr];
      const source = previous && current === previous.translated ? previous.source : current;
      const translated = translateText(source);
      values[attr] = { source, translated };
      if (current !== translated) element.setAttribute(attr, translated);
    }
    element.childNodes.forEach(translateElement);
  }

  function updateSwitcher(switcher) {
    if (switcher.classList.contains('site-language-toggle')) {
      switcher.setAttribute('aria-label', language === 'en' ? 'Language' : 'Ngôn ngữ');
      switcher.querySelectorAll('[data-language]').forEach((button) => {
        const selected = button.dataset.language === language;
        button.setAttribute('aria-pressed', String(selected));
        button.setAttribute('aria-label', button.dataset.language === 'en' ? 'English' : 'Tiếng Việt');
      });
      return;
    }
    switcher.value = language;
    switcher.setAttribute('aria-label', language === 'en' ? 'Language' : 'Ngôn ngữ');
    switcher.title = language === 'en' ? 'Language' : 'Ngôn ngữ';
  }

  function setLanguage(nextLanguage) {
    language = nextLanguage === 'en' ? 'en' : 'vi';
    localStorage.setItem(STORAGE_KEY, language);
    document.documentElement.lang = language;
    document.querySelectorAll('input, select, textarea').forEach((control) => {
      control.setCustomValidity('');
    });
    applying = true;
    translateElement(document.body);
    document.querySelectorAll('.site-language-switcher').forEach(updateSwitcher);
    if (document.title) {
      const titleNode = document.querySelector('title');
      if (titleNode) translateElement(titleNode);
    }
    applying = false;
  }

  function getValidationMessage(control) {
    if (control.validity.valueMissing) {
      return language === 'en' ? 'This field is required.' : 'Vui lòng điền trường này.';
    }
    if (control.validity.typeMismatch && control.type === 'email') {
      return language === 'en' ? 'Please enter a valid email address.' : 'Vui lòng nhập địa chỉ email hợp lệ.';
    }
    if (control.validity.tooShort) {
      const length = control.minLength;
      return language === 'en'
        ? `Please use at least ${length} characters.`
        : `Vui lòng nhập ít nhất ${length} ký tự.`;
    }
    if (control.validity.tooLong) {
      const length = control.maxLength;
      return language === 'en'
        ? `Please use no more than ${length} characters.`
        : `Vui lòng nhập không quá ${length} ký tự.`;
    }
    return language === 'en' ? 'Please check this field.' : 'Vui lòng kiểm tra trường này.';
  }

  document.addEventListener('invalid', (event) => {
    const control = event.target;
    if (control instanceof HTMLInputElement || control instanceof HTMLSelectElement || control instanceof HTMLTextAreaElement) {
      control.setCustomValidity(getValidationMessage(control));
    }
  }, true);
  document.addEventListener('input', (event) => {
    if (event.target instanceof HTMLInputElement || event.target instanceof HTMLTextAreaElement) {
      event.target.setCustomValidity('');
    }
  }, true);
  document.addEventListener('change', (event) => {
    if (event.target instanceof HTMLSelectElement) event.target.setCustomValidity('');
  }, true);

  function createSwitcher(segmented = false) {
    if (segmented) {
      const group = document.createElement('div');
      group.className = 'site-language-switcher site-language-toggle';
      group.setAttribute('role', 'group');
      group.setAttribute('aria-label', 'Ngôn ngữ');
      for (const code of ['vi', 'en']) {
        const button = document.createElement('button');
        button.type = 'button';
        button.className = 'language-option';
        button.dataset.language = code;
        button.textContent = code.toUpperCase();
        button.setAttribute('aria-label', code === 'en' ? 'English' : 'Tiếng Việt');
        button.setAttribute('aria-pressed', String(code === language));
        button.addEventListener('click', () => setLanguage(code));
        group.appendChild(button);
      }
      return group;
    }
    const select = document.createElement('select');
    select.className = 'site-language-switcher';
    select.setAttribute('aria-label', 'Language');
    select.title = 'Language';
    select.innerHTML = '<option value="vi">Tiếng Việt</option><option value="en">English</option>';
    select.value = language;
    select.addEventListener('change', () => setLanguage(select.value));
    return select;
  }

  function mountSwitcher(container) {
    if (!container || container.querySelector('.site-language-switcher')) return;
    const floatingSwitcher = document.querySelector('.site-language-switcher-floating');
    if (floatingSwitcher) {
      floatingSwitcher.classList.remove('site-language-switcher-floating');
      container.appendChild(floatingSwitcher);
    } else {
      container.appendChild(createSwitcher(
        container.classList.contains('auth-language-control') ||
        container.classList.contains('login-language-control')
      ));
    }
  }

  function initialize() {
    document.documentElement.lang = language;
    const loginLanguage = document.querySelector('.auth-language-control, .login-language-control');
    const navRight = document.querySelector('.navbar-right');
    if (loginLanguage) {
      mountSwitcher(loginLanguage);
    } else if (navRight) {
      mountSwitcher(navRight);
    } else {
      const switcher = createSwitcher();
      switcher.classList.add('site-language-switcher-floating');
      document.body.appendChild(switcher);
    }
    setLanguage(language);

    new MutationObserver((changes) => {
      if (applying) return;
      changes.forEach((change) => {
        if (change.type === 'characterData') translateNode(change.target);
        if (change.type === 'attributes') translateElement(change.target);
        change.addedNodes?.forEach(translateElement);
      });
    }).observe(document.body, {
      childList: true,
      characterData: true,
      attributes: true,
      attributeFilter: ['placeholder', 'title', 'aria-label', 'alt', 'value'],
      subtree: true
    });

    new MutationObserver(() => {
      const navRightNow = document.querySelector('.navbar-right');
      if (navRightNow) mountSwitcher(navRightNow);
    }).observe(document.body, { childList: true, subtree: true });
  }

  window.SiteI18n = { get language() { return language; }, setLanguage, mountSwitcher, translate: translateText };
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initialize, { once: true });
  } else {
    initialize();
  }
})();
