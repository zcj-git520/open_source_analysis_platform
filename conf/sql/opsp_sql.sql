INSERT INTO `language` (
    `name`,
    `image`,
    `score`,
    `desc`,
    `repo_rul`,
    `bio`
) VALUES
(
    'Python',
    'https://upload.wikimedia.org/wikipedia/commons/thumb/c/c3/Python-logo-notext.svg/1200px-Python-logo-notext.svg.png',
    '5.9', -- 假设的TIOBE分数
    'Python是一种广泛使用的高级编程语言，以其简洁易读的语法和强大的库支持而闻名。',
    'https://github.com/python/cpython',
    'Python由Guido van Rossum于1991年首次发布，是数据科学、人工智能、网络开发等领域的首选语言之一。'
),
(
    'JavaScript',
    'https://upload.wikimedia.org/wikipedia/commons/thumb/9/99/Unofficial_JavaScript_logo_2017.svg/1200px-Unofficial_JavaScript_logo_2017.svg.png',
    '5.5', -- 假设的TIOBE分数
    'JavaScript是网页开发的核心技术之一，也广泛用于服务器端编程和移动应用开发。',
    'https://github.com/tc39/ecma262',
    'JavaScript由Netscape公司在1995年首次推出，如今已成为Web开发不可或缺的一部分。'
),
(
    'Java',
    'https://upload.wikimedia.org/wikipedia/commons/thumb/d/de/Java_Programming_Language_logo.svg/1200px-Java_Programming_Language_logo.svg.png',
    '4.8', -- 假设的TIOBE分数
    'Java是一种面向对象的编程语言，广泛应用于企业级应用、Android应用开发、大数据处理等领域。',
    'https://github.com/openjdk',
    'Java拥有庞大的生态系统，是许多大型企业和组织的首选语言。'
),
(
    'C++',
    'https://upload.wikimedia.org/wikipedia/commons/thumb/1/18/ISO_C%2B%2B_Logo.svg/1200px-ISO_C%2B%2B_Logo.svg.png',
    '4.5', -- 假设的TIOBE分数
    'C++是一种高效、灵活的编程语言，支持面向对象编程、泛型编程和过程化编程。',
    'https://github.com/cplusplus/draft',
    'C++由Bjarne Stroustrup于1980年代初开发，是系统/应用软件开发、游戏开发、嵌入式系统等领域的重要语言。'
),
(
    'C#',
    'https://upload.wikimedia.org/wikipedia/commons/thumb/4/4c/CSharp_Logo.svg/1200px-CSharp_Logo.svg.png',
    '4.0', -- 假设的TIOBE分数
    'C#是一种面向对象的编程语言，由微软开发，是.NET框架的核心语言。',
    'https://github.com/dotnet/csharplang',
    'C#广泛用于Windows应用开发、游戏开发、Web开发等领域，是微软生态系统中的重要组成部分。'
),
(
    'PHP',
    'https://upload.wikimedia.org/wikipedia/commons/thumb/2/27/PHP-logo.svg/1200px-PHP-logo.svg.png',
    '3.9', -- 假设的TIOBE分数
    'PHP是一种广泛使用的开源通用脚本语言，特别适合于Web开发，能够嵌入到HTML中。',
    'https://github.com/php/php-src',
    'PHP由Rasmus Lerdorf在1994年创建，是许多动态网站和Web应用程序的后端语言。'
),
(
    'C',
    'https://upload.wikimedia.org/wikipedia/commons/thumb/3/35/The_C_programming_language_logo.svg/1200px-The_C_programming_language_logo.svg.png',
    '10.2', -- 假设的 TIOBE 分数，可根据实际情况调整
    'C 语言是一种广泛使用的通用编程语言，具有高效、灵活和可移植性等特点。它是许多其他编程语言的基础，被广泛应用于系统软件、嵌入式系统、游戏开发等领域。',
    'https://github.com/curl/curl', -- 这里只是一个 C 语言相关项目的仓库地址示例，可根据实际情况调整
    'C 语言以其简洁的语法和强大的功能而闻名，是编程领域的经典语言之一。'
),
(
    'Ruby',
    'https://upload.wikimedia.org/wikipedia/commons/thumb/c/c3/Ruby_logo.svg/1200px-Ruby_logo.svg.png',
    '2.9', -- 假设的TIOBE分数
    'Ruby是一种动态、反射式、面向对象的编程语言，以其简洁和强大的功能而著称。',
    'https://github.com/ruby/ruby',
    'Ruby由松本行弘（Yukihiro Matsumoto）设计，是Web开发、脚本编写和快速原型设计的理想选择。'
),
(
    'Golang',
    'https://upload.wikimedia.org/wikipedia/commons/thumb/c/c7/Go_Logo_Blue.svg/1200px-Go_Logo_Blue.svg.png',
    '2.5', -- 假设的TIOBE分数
    'Go（又称Golang）是一种由Google开发的静态类型、编译型编程语言，设计用于构建简单、可靠和高效的软件。',
    'https://github.com/golang/go',
    'Go语言以其简洁的语法、高效的并发支持和强大的标准库而受到开发者的青睐，广泛用于云计算、微服务等领域。'
),

(
    'Rust',
    'https://upload.wikimedia.org/wikipedia/commons/thumb/0/0c/Rust_programming_language_black_logo.svg/1200px-Rust_programming_language_black_logo.svg.png',
    '2.3', -- 假设的TIOBE分数或其他排名指标
    'Rust是一种注重安全、性能和并发的系统编程语言，由Mozilla主导开发。',
    'https://github.com/rust-lang/rust',
    'Rust以其内存安全性和对并发编程的支持而受到赞誉，被广泛应用于系统编程、游戏开发、网络服务和嵌入式系统等领域。'
),
(
    'TypeScript',
    'https://upload.wikimedia.org/wikipedia/commons/thumb/2/2c/TypeScript_logo_2020.svg/1200px-TypeScript_logo_2020.svg.png',
    '2.8', -- 假设的流行度分数
    'TypeScript是JavaScript的一个超集，添加了可选的静态类型和基于类的面向对象编程。',
    'https://github.com/microsoft/TypeScript',
    'TypeScript由Microsoft开发，旨在解决JavaScript在大型应用中的可维护性和可扩展性问题，广泛应用于前端和Node.js开发。'
);

-- Web 前端开发
INSERT INTO `repo_category` (`name`, `image_url`, `desc`)
VALUES
    ("Web前端开发", "https://10.idqqimg.com/qqcourse_logo_ng/ajNVdqHZLLCUfI85kcsUzqwyPLb2zLM4loia5rqMFHp07bYOOXWttZrvyCfJKfMmgnQ2wyYYoXvA/600", "前端框架如:React、Vue.js、Angular等现代前端框架的项目非常热门"),
    ("Web后端框架", "https://atts.w3cschool.cn/attachments/image/20201208/1607410145715356.jpg", "包括后端框架（如 Django、Flask）等"),
    ("移动应用开发", "https://pic.pngsucai.com/00/26/94/920fc444f58f6b18.webp", "例如 Android 和 iOS 应用开发项目"),
    ("数据科学与机器学习", "https://swarma.org/wp-content/uploads/2020/07/wxsync-2020-07-29a86936de490cdc282ffba1fd92cb56.png", "包含数据分析工具、机器学习算法实现等"),
    ("游戏开发", "https://img2020.cnblogs.com/blog/1561139/202010/1561139-20201012165557046-933664682.png", "涵盖各种游戏引擎和游戏项目"),
    ("系统工具", "https://www.zmtests.com/image/8da0447b-c397-4c8b-b808-9979dec1cf36", "如操作系统相关项目、命令行工具等"),
    ("开源库和框架", "https://tse2-mm.cn.bing.net/th/id/OIP-C.AZOMg3sc34S3k_LEAkJD2QHaDt?rs=1&pid=ImgDetMain", "为其他开发者提供功能模块的项目，如数据库连接库、图形处理库等"),
    ("人工智能与机器学习", "https://pic4.zhimg.com/v2-34d43321006cffe0b5ce509a39913959_r.jpg?source=172ae18b", "包括深度学习模型、自然语言处理等项目"),
    ("区块链", "https://pic4.zhimg.com/v2-da0b06a82473d898457d8b7f77473026_r.jpg", "涉及区块链技术的实现和应用"),
    ("物联网（IoT）", "https://www.shuangyi-tech.com/upload/month_2011/202011261009529139.jpg", "与物联网设备和系统相关的项目"),
    ("教育资源", "https://pic3.zhimg.com/v2-69d8f7909d139258d04d2f16ad6477e0_1440w.jpg?source=172ae18b", "如编程教程、学习项目等"),
    ("安全与加密", "https://tse3-mm.cn.bing.net/th/id/OIP-C.0L5yfxYdb3qZnZJzRBWoPQHaDJ?rs=1&pid=ImgDetMain", "包括安全工具、加密算法实现等"),
    ("设计与创意", "https://img.zcool.cn/community/01b85255f3e60a6ac7251df8682e78.jpg@1280w_1l_2o_100sh.jpg", "如图形设计工具、创意代码项目等"),
    ("其他","https://www.shejizn.vip/wp-content/uploads/2020/07/1594349552-8f14e45fceea167.png","其他");
