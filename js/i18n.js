(function () {
  const STORAGE_KEY = 'vinhomes-language';
  const translations = {
    'BẢNG ĐIỀU KHIỂN GIS / BIM 3D': '3D GIS / BIM CONTROL PANEL',
    '01 CẮT LỚP CHIỀU CAO & THI CÔNG': '01 HEIGHT SECTION & CONSTRUCTION',
    '02 QUẢN LÝ LỚP DỮ LIỆU': '02 DATA LAYER MANAGEMENT',
    '03 THÔNG TIN PHÂN KHU CHỌN': '03 SELECTED AREA INFORMATION',
    'CẮT LỚP CHIỀU CAO & THI CÔNG': 'HEIGHT SECTION & CONSTRUCTION',
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
    'Cư dân': 'Residents',
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
    'Tiếng Việt': 'Vietnamese'
  };
  const entries = Object.entries(translations).sort((a, b) => b[0].length - a[0].length);
  const originals = new WeakMap();
  let language = localStorage.getItem(STORAGE_KEY) === 'en' ? 'en' : 'vi';
  let applying = false;

  function translateText(value) {
    if (language !== 'en') return value;
    return entries.reduce((text, [source, target]) => text.split(source).join(target), value);
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

    for (const attr of ['placeholder', 'title', 'aria-label']) {
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
    switcher.value = language;
    switcher.setAttribute('aria-label', language === 'en' ? 'Language' : 'Ngôn ngữ');
    switcher.title = language === 'en' ? 'Language' : 'Ngôn ngữ';
  }

  function setLanguage(nextLanguage) {
    language = nextLanguage === 'en' ? 'en' : 'vi';
    localStorage.setItem(STORAGE_KEY, language);
    document.documentElement.lang = language;
    applying = true;
    translateElement(document.body);
    document.querySelectorAll('.site-language-switcher').forEach(updateSwitcher);
    if (document.title) {
      const titleNode = document.querySelector('title');
      if (titleNode) translateElement(titleNode);
    }
    applying = false;
  }

  function createSwitcher() {
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
      container.appendChild(createSwitcher());
    }
  }

  function initialize() {
    document.documentElement.lang = language;
    const navRight = document.querySelector('.navbar-right');
    if (navRight) {
      mountSwitcher(navRight);
    } else {
      const switcher = createSwitcher();
      switcher.classList.add('site-language-switcher-floating');
      document.body.appendChild(switcher);
    }
    setLanguage(language);

    new MutationObserver((changes) => {
      if (applying || language !== 'en') return;
      changes.forEach((change) => {
        if (change.type === 'characterData') translateNode(change.target);
        change.addedNodes?.forEach(translateElement);
      });
    }).observe(document.body, { childList: true, characterData: true, subtree: true });

    new MutationObserver(() => {
      const navRightNow = document.querySelector('.navbar-right');
      if (navRightNow) mountSwitcher(navRightNow);
    }).observe(document.body, { childList: true, subtree: true });
  }

  window.SiteI18n = { get language() { return language; }, setLanguage, mountSwitcher };
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initialize, { once: true });
  } else {
    initialize();
  }
})();
