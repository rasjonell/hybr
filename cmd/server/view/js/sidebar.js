function toggleSidebar() {
  const sidebar = document.getElementById('sidebar');
  const main = document.getElementById('main-content');
  const expandIcon = document.getElementById('expand-icon');
  const collapseIcon = document.getElementById('collapse-icon');
  const isExpanded = sidebar.classList.contains('w-64');

  sidebar.classList.toggle('w-64');
  sidebar.classList.toggle('w-16');

  main.classList.toggle('ml-64');
  main.classList.toggle('ml-16');

  console.log('expand hidden?', expandIcon.classList.contains('hidden'))
  console.log('collapse hidden?', collapseIcon.classList.contains('hidden'))

  expandIcon.classList.toggle('hidden');
  collapseIcon.classList.toggle('hidden');


  const textElements = document.querySelectorAll('.sidebar-text');
  textElements.forEach(el => {
    if (isExpanded) {
      el.classList.add('hidden');
      el.classList.remove('inline', 'block');
    } else {
      el.classList.remove('hidden');
      if (el.tagName.toLowerCase() === 'p') {
        el.classList.add('block');
      } else {
        el.classList.add('inline');
      }
    }
  });

  const tooltips = sidebar.querySelectorAll('.group-hover\\:block');
  tooltips.forEach(el => {
    if (isExpanded) {
      el.classList.add('group-hover:block');
      el.classList.remove('group-hover:hidden');
    } else {
      el.classList.add('group-hover:hidden');
      el.classList.remove('group-hover:block');
    }
  });
}

document.addEventListener('DOMContentLoaded', function() {
  const isLargeScreen = window.matchMedia('(min-width: 1024px)').matches;

  if (isLargeScreen) {
    toggleSidebar();
  }
});

window.addEventListener('resize', function() {
  const sidebar = document.getElementById('sidebar');
  const isLargeScreen = window.matchMedia('(min-width: 1024px)').matches;
  const isExpanded = sidebar.classList.contains('w-64');

  if (isLargeScreen && !isExpanded) {
    toggleSidebar();
  } else if (!isLargeScreen && isExpanded) {
    toggleSidebar();
  }
});

