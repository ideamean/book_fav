var bookFav = {
	getBookList: function () {
		return 123;
	}
};
var app = new Vue({
	el: '#app',
	data: {
		app: {
			title: "家庭图书馆",
			userInfo: {}
		},
		hasMore: 1,
		getMoreTxt: '查看更多',
		showSearch: false, //默认不展示搜索框
		isbn: "",
		showLogin: false,
		userEmail: '',
		genCode: '',
		genCodeMessage: '',
		bookData: {},
		searchTipMessage: '',
		searchResultData: null,
		purchaseRemark: '', //加入已购备注
		bookListType: 1
 	},
	methods: {
		init: function () {
			document.title = this.app.title;
			this.getUserInfo();
		},
		getUserInfo: function () {
			var o = this;
			axios.get('/api/user/info').then(function (response) {
				var data = response.data;
				if (data.errno > 0) {
					if (data.errno == 100000) {
						o.showLogin = true
					}
					return;
				}
				o.app.userInfo = data.data;
				o.getBookList(1, 20);
			}).catch(function (error) {
				console.log(error);
			});
		},
		login: function () {
			var o = this;
			o.genCodeMessage = '';
			var email = o.userEmail;
			var genCode = o.genCode;
			axios.get('/api/user/login?email='+email+'&log_type=login&gen_code=' + genCode).then(function (response) {
				var data = response.data;
				if (data.errno > 0) {
					o.genCodeMessage = data.errmsg;
					return;
				}
				o.showLogin = false
				o.init();
			}).catch(function (error) {
				o.genCodeMessage = error;
			});
		},
		logout: function () {
			var o = this;
			axios.get('/api/user/logout').then(function (response) {
				o.showLogin = true
			}).catch(function (error) {
				console.log(error);
			});
		},
		sendEmailCode: function () {
			var email = this.userEmail;
			if (email == '') {
				this.genCodeMessage = '请输入邮箱地址';
				return;
			}

			var o = this;
			axios.get('/api/user/login?email=' + email + '&log_type=send_code').then(function (response) {
				var data = response.data;
				if (data.errno > 0) {
					o.genCodeMessage = data.errmsg;
					return;
				}
				o.genCodeMessage = '已发送,请打开邮箱查看验证码';
			}).catch(function (error) {
				o.genCodeMessage = error;
			});
		},
		bookSearchByIsbn: function () {
			var isbn = this.isbn;
			var o = this;
			axios.get('/api/book/search?isbn=' + isbn).then(function (response) {
				var data = response.data;
				if (data.errno > 0) {
					o.searchTipMessage = data.errmsg;
					return;
				}
				o.searchTipMessage = '';
				o.searchResultData = data.data;
			}).catch(function (error) {
				o.searchTipMessage = error;
			});
		},
		addPurchase: function (bid) {
			var remark = this.purchaseRemark;
			var o = this;
			axios.get('/api/book/addpurchase?bid=' + bid + '&remark=' + remark).then(function (response) {
				var data = response.data;
				if (data.errno > 0) {
					o.searchTipMessage = data.errmsg;
					return;
				}
				o.searchTipMessage = '';
				o.searchResultData.purchase = {is_purchase: 2};
			}).catch(function (error) {
				o.searchTipMessage = error;
			});
		},
		addPurchaseFromList: function (index) {
			var bookInfo = this.bookData.book_list[index];
			if (bookInfo.isPurchase) {
				return;
			}
			this.searchResultData = {
				purchase: {
					is_purchase: 1,
				},
				book_info: bookInfo
			};
			document.body.scrollTop = 0
			document.documentElement.scrollTop = 0
		},
		setBookListType: function (type) {
			this.bookListType = type;
			this.hasMore = 0;
			this.getMoreTxt = '查看更多';
			this.getBookList(1, 20);
		},
		getBookList: function (page, size) {
			var o = this;
			var type = this.bookListType;
			axios.get('/api/book/list?type=' + type + '&page=' + page + '&size=' + size).then(function (response) {
			    var data = response.data;
				if (data.errno > 0) {
					alert(data.errmsg);
					return;
				}
				var s = data.data;
				if (s.page <= 1) {
					o.bookData = s;
					o.hasMore = 1;
				}

				o.bookData.page = page;
				o.bookData.page_size = size;

				if (s.book_list == null) {
					var len = 0;
				} else {
					var len = s.book_list.length;
				}

				for (var i = 0; i < len; i++) {
					var bid = s['book_list'][i].bid;
					s['book_list'][i].isPurchase = 0;
					if (s['purchase_list'][bid]) {
						s['book_list'][i].isPurchase = 1;
					}
					if (page > 1) {
						o.bookData.book_list.push(s['book_list'][i]);
					}
				}

				if (len <= 0) {
					o.hasMore = 0;
					o.getMoreTxt = '到底啦';
				}

			}).catch(function (error) {
				console.log(error);
			});
		},
		getMore: function () {
			var nextPage = this.bookData.page - 0 + 1;
			this.getBookList(nextPage, this.bookData.page_size);
		},
		showSearchBar: function () {
			this.showSearch = !this.showSearch;
		},
		imageUploadBefore: function () {
			this.searchTipMessage = '正在上传...';
			return true;
		},
        imageUploadSuccess: function (response, file, fileList) {
        	this.searchTipMessage = '';
            if (response.errno > 0) {
            	this.searchTipMessage = response.errmsg;
            	return;
            }
            var isbnArr = response.data;
            if (isbnArr == null) {
            	this.searchTipMessage = "示识别出条形码!";
            	return;
            }
            this.isbn = isbnArr[0];
            this.bookSearchByIsbn();
        }
	}
});

app.init();




