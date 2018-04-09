$(document).ready(function() {

  $(".baner__close").click(function() {
    $('.baner').fadeOut(100);  
  });

});
jQuery(document).ready(function($) {
  $(".doc-item__name").dotdotdot({
    height: 43
  });

  if (window.matchMedia("(max-width: 767px)").matches) {
    $(".doc-item__more").click(function() {
      $(".doc-item__name").dotdotdot({
        height: 60
      });
    });
  }

  if (window.matchMedia("(min-width: 768px)").matches) {
    // Высота doc-item__more-block = высоте doc-item__header
    $(".doc-item__more-block").css("top", parseInt($(".doc-item__header").innerHeight() - 5));
  }

  // fadeToggle для doc-item__more-block
  $(".doc-item__more").click(function() {
    if (window.matchMedia("(min-width: 768px)").matches) {
      $(this).parents(".doc-item").find(".doc-item__more-block").fadeToggle();
    }
  });

  // Высота блока doc-item__img-block--add
  var heightDocItem = parseInt($(".doc-item__img-block").height());

  $(".doc-item__img-block--add").height(heightDocItem);
  docItemHover();

  $(window).resize(function() {
    once(function() {
      var heightDocItem = parseInt($(".doc-item__img-block").height());
      $(".doc-item__img-block--add").height(heightDocItem);

      // Высота doc-item__more-block = высоте doc-item__header
      $(".doc-item__more-block").css("top", parseInt($(".doc-item__header").innerHeight() - 5));

      // Ховер на doc-item
      docItemHover();

      if (window.matchMedia("(min-width: 768px)").matches) {
        $(".doc-item__name").dotdotdot({
          height: 43
        });
      } else if (window.matchMedia("(max-width: 767px)").matches) {
        $(".doc-item").unbind("mouseenter mouseleave");
      }
    });
  });
});

// Ховер на doc-item
function docItemHover() {
  $(".doc-item").hover(function() {
    if (window.matchMedia("(min-width: 768px)").matches) {
      var itemHeight = $(this).innerHeight();
      var headerHeight = $(this).find(".doc-item__header").innerHeight();
      $(this).find(".doc-item__content").css("minHeight", itemHeight - headerHeight);

      $(this).css("minHeight", $(this).height());

      var initianHeight = parseInt($(this).find(".doc-item__content").height());
      $(this).find(".doc-item__name").dotdotdot({
        height: null
      });
      $(this).addClass("doc-item--hover");

      var finalHieght = parseInt($(this).find(".doc-item__content").height());
      $(this).find(".doc-item__content").css("bottom", initianHeight - finalHieght);
      $(this).find(".doc-item__content").addClass("doc-item__content--hover");
    }
  },
  function() {
    if (window.matchMedia("(min-width: 768px)").matches) {
      $(this).find(".doc-item__name").dotdotdot({
        height: 43
      });
      $(this).removeAttr("style");

      $(this).removeClass("doc-item--hover");
      $(this).find(".doc-item__content").removeClass("doc-item__content--hover");
      $(this).find(".doc-item__more-block").fadeOut();
      $(this).find(".doc-item__content").css("bottom", "");
    }
  });
}

// чтобы функция выполнилась только один раз https://frontender.info/essential-javascript-functions/
function once(fn, context) { 
  var result;

  return function() { 
    if (fn) {
      result = fn.apply(context || this, arguments);
      fn = null;
    }

    return result;
  };
}
$(document).ready(function () {
	console.log($(".swiper-container").width())
	var slide = 3;
	var center = true;
	if ($(".swiper-container").width() < 567) {
		slide = 2;
		center = false;
	}
	if ($(".swiper-container").width() < 362) {
		slide = 1;
		center = true;
	}
	//initialize swiper when document ready
	var mySwiper = new Swiper ('.swiper-container', {	  
	  loop: true,
	  slidesPerView: slide,
      spaceBetween: 50,
      centeredSlides: center,
	  navigation: {
	    nextEl: '.doc-templates__slider-next',
	    prevEl: '.doc-templates__slider-prev',
	  },
	})
});
$(".form__button--add-img").on("change", function() {
  var files = this.files;

  for (var i = 0; i < files.length; i++) {
    var file = files[i];

    if (!file.type.match(/image\/(jpeg|jpg|png|gif)/)) {
      alert("Фотография должна быть в формате jpg, png или gif");
      continue;
    }

    if (file.size > maxFileSize) {
      alert("Размер фотографии не должен превышать 2 Мб");
      continue;
    }

    preview(files[i]);
  }

  this.value = "";
});

$(".input-file__input").on("change", function() {
  var files = this.files;
  for (var i = 0; i < files.length; i++) {
    var file = files[i];

    if (!file.type.match(/image\/(jpeg|jpg|png|gif)/)) {
      alert("Фотография должна быть в формате jpg, png или gif");
      continue;
    }

    fileName = file.name.replace(/\\/g, "/").split("/").pop();
    $(this).parents(".input-file").find(".input-file__file-name").text(fileName);
  }
});
var text = $('.iframe__description').text()
var symbol= text.length;
console.log(text)
if (symbol > 400) {
	var newText;
	for (var i = 0; i < 400; i++) {
		newText = newText + text[i];
	}
	$('.iframe__description').text(newText);
	$('.iframe__description').addClass('iframe__description--after');
	console.log(newText);
}
$( "#iframe__button--open-text" ).click(function() {
  $('.iframe__description').text(text);
  $('.iframe__description').removeClass('iframe__description--after');
});


$(".lang").click(function() {
  $(".lang__drop-box").toggleClass("lang__drop-box--open");
});
function equalWidthListNum(numClass) {
  var listNumMaxWidth = 0;

  $(numClass).each(function() {
    if ($(this).width() > listNumMaxWidth) {
      listNumMaxWidth += $(this).width();
    }
  });

  $(numClass).css("width", listNumMaxWidth);
}
$(document).ready(function() {
  equalWidthListNum(".list__item-num");

  $(".list__title").click(function() {
    $(this).siblings(".list--content").slideToggle();
  });

  $(".list__item--miltilevel .list__icon--right-arrow").click(function(e) {
    $(this).parent().siblings(".list--miltilevel").slideToggle();
  });
});
$(document).ready(function() {
  $(".nav__toggle").click(function() {
    $(".header__right-side").fadeToggle(function() {
      if ($(this).is(":visible")) {
        $(this).css("display", "flex");
      }
    });
    $(".page").toggleClass("page--scroll");
  });
});
svg4everybody();
$(".popup-modal__doc-item").magnificPopup({
  type: "inline",
  removalDelay: 400,
  mainClass: "my-mfp-zoom-in",
  fixedContentPos: true,
  disableOn: function() {
    if (window.matchMedia("(max-width: 767px)").matches) {
      return true;
    } else if (window.matchMedia("(min-width: 768px)").matches) {
      return false;
    }

    $(window).resize(function() {
      if (window.matchMedia("(max-width: 767px)").matches) {
        return true;
      } else if (window.matchMedia("(min-width: 768px)").matches) {
        return false;
      }
    });
  }
});

$(".popup-modal").magnificPopup({
  type: "inline",
  removalDelay: 400,
  mainClass: "my-mfp-zoom-in",
  fixedContentPos: true
});

$(".form__button--popup-ok").click(function() {
  $.magnificPopup.close();
});

$(".form__button--send").click(function() {
  blockCenter(".popup--ok");
});

$(window).resize(function() {
  blockCenter(".popup--ok");
});

// Замена translate(-50%, -50%)
function blockCenter(block) {
  var halfWidthWindow = $(window).width() / 2;
  var halfHeightWindow = $(window).height() / 2;
  $(block).css({
    "top": halfHeightWindow - $(block).innerHeight() / "2",
    "left": halfWidthWindow - $(block).innerWidth() / "2"
  });
}
$(document).ready(function() {
  if (window.matchMedia("(min-width: 992px)").matches) {
    $(".sidebar").niceScroll({
      cursorcolor: "#eaeaea",
      cursorwidth: "4px",
      background: "#f7f7f7",
      nativeparentscrolling: false
    });

    /* if ($(window).scrollTop() > 54) {
      $(".sidebar").css("top", 0);
    } else {
      $(".sidebar").css("top", 54);
    } */

    /* attachesSidebarTop(); */
  }

  $(".sidebar__toggle").click(function() {
    if ($(this).find("use").attr("xlink:href") === "imgs/sprite.svg#delete") {
      $(this).find("use").attr("xlink:href", "imgs/sprite.svg#sidebar");
      $(this).removeClass("sidebar__toggle--delete");
    } else if ($(this).find("use").attr("xlink:href") === "imgs/sprite.svg#sidebar") {
      $(this).find("use").attr("xlink:href", "imgs/sprite.svg#delete");
      $(this).addClass("sidebar__toggle--delete");
    }
    $(".sidebar").fadeToggle();
  });

  $(".form__button--help").click(function() {
    $(".sidebar").fadeToggle();
  });

  $(window).resize(function() {
    if (window.matchMedia("(min-width: 992px)").matches) {
      $(".sidebar").show();
    }
  });

  $(".user").click(function() {
    $(this).parents(".sidebar__user-container").find(".list--company-list").slideToggle();
  });
});

/* function attachesSidebarTop() {
  var currTop = 54;
  var currPaddingBot = parseInt($(".sidebar").css("paddingBottom"));

  $(window).scroll(function() {
    $(".sidebar").css({"top": function() {
      if ($(window).scrollTop() < 54) {
        return currTop - $(window).scrollTop();
      } else {
        return 0;
      }
    },
    "paddingBottom": function() {
      if ($(window).scrollTop() < 54) {
        return currPaddingBot - $(window).scrollTop();
      } else {
        return 20;
      }
    }
    });
  });
} */

$(document).ready(function() {
  $(".sidebar-folder__content-toggle").click(function() {
    $(this).parents(".sidebar-folder").toggleClass("sidebar-folder--open");
    $(this).parents(".sidebar-folder").find(".sidebar-folder__content").slideToggle();
  });
});
function openTab(evt, tabName, speed) {
  let container = ".tabs";
  let buttonsContainer = ".tabs__nav";
  let buttons = ".tabs__button";
  let content = ".tabs__content";

  // Проверка на наличие индентификатора.
  if ($(evt.target).parents(container).find(content).is("#" + tabName)) {
    // evt.target - таб, на который был клик.
    $(evt.target).parent(buttonsContainer).find(buttons).removeClass("tabs__button--active");
    $(evt.target).parents(container).find(content).hide();

    $("#" + tabName).fadeIn(speed * 2);
    $(evt.target).addClass("tabs__button--active");

    $(".doc-item__name").dotdotdot({
      height: 43
    });
  } else {
    alert("Неверный идентификатор!");
  }
}

// Открывает первый таб.
$(".tabs").each(function() {
  $(this).find(".tabs__button:first").trigger("click");
});