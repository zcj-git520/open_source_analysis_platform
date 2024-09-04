INSERT INTO `osap_user`.`language` (
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
    'C/C++',
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
    'Swift',
    'https://upload.wikimedia.org/wikipedia/commons/thumb/a/a7/Swift_logo.svg/1200px-Swift_logo.svg.png',
    '3.5', -- 假设的TIOBE分数
    'Swift是一种由苹果公司开发的编程语言，用于iOS、macOS、watchOS和tvOS应用开发。',
    'https://github.com/apple/swift',
    'Swift以其简洁、快速和安全的特性而受到开发者的喜爱，是苹果生态系统中的主要编程语言之一。'
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
    'Go',
    'https://upload.wikimedia.org/wikipedia/commons/thumb/c/c7/Go_Logo_Blue.svg/1200px-Go_Logo_Blue.svg.png',
    '2.5', -- 假设的TIOBE分数
    'Go（又称Golang）是一种由Google开发的静态类型、编译型编程语言，设计用于构建简单、可靠和高效的软件。',
    'https://github.com/golang/go',
    'Go语言以其简洁的语法、高效的并发支持和强大的标准库而受到开发者的青睐，广泛用于云计算、微服务等领域。'
),
(
    'Kotlin',
    'https://upload.wikimedia.org/wikipedia/commons/thumb/a/a6/Kotlin_logo.svg/1200px-Kotlin_logo.svg.png',
    '2.0', -- 假设的TIOBE分数
    'Kotlin是一种静态类型编程语言，由JetBrains设计并开源，用于现代多平台应用。',
    'https://github.com/JetBrains/kotlin',
    'Kotlin与Java高度互操作，是Android应用开发的首选语言之一，也广泛用于服务器端和Web开发。'
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
),
(
    'Visual Basic .NET (VB.NET)',
    'https://upload.wikimedia.org/wikipedia/commons/thumb/3/35/Visual_Basic_.NET_Logo.svg/1200px-Visual_Basic_.NET_Logo.svg.png',
    '1.5', -- 假设的流行度分数
    'Visual Basic .NET（VB.NET）是一种面向对象的编程语言，是.NET框架的一部分。',
    '（注意：VB.NET可能没有像其他语言那样明显的官方GitHub仓库，这里不提供URL）',
    'VB.NET是Visual Basic的继承者，旨在与.NET框架一起使用，提供快速的应用程序开发（RAD）能力，尤其适合于Windows应用开发。'
),
(
    'R',
    'https://upload.wikimedia.org/wikipedia/commons/thumb/1/1a/R_logo.svg/1200px-R_logo.svg.png',
    '1.8', -- 假设的流行度分数
    'R是一种用于统计计算和图形的编程语言和软件环境。',
    'https://github.com/wch/r-source',
    'R以其强大的统计分析和图形功能而受到数据科学家、统计学家和经济学家的喜爱，广泛用于学术研究、商业分析和数据挖掘。'
),
(
    'Scala',
    'https://upload.wikimedia.org/wikipedia/commons/thumb/a/a2/Scala_logo.svg/1200px-Scala_logo.svg.png',
    '1.2', -- 假设的流行度分数
    'Scala是一种多范式编程语言，旨在以简洁、优雅和类型安全的方式表达常见的编程模式。',
    'https://github.com/scala/scala',
    'Scala运行在Java虚拟机（JVM）上，与Java高度互操作，被广泛应用于大数据处理、Web开发、机器学习等领域。'
);
-- 你可以继续添加其他热门语言的信息
;