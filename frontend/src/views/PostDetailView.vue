<template>
  <article v-if="post" class="fixed inset-0 h-screen overflow-hidden flex flex-col bg-background"
           :class="post.postType === 'article' ? 'pt-24' : 'pt-8'">
    <!-- Header: Fixed at top (Hidden on desktop for gallery) -->
    <header
        class="transition-opacity duration-300 space-y-3 max-w-[1024px] mx-auto px-8 w-full shrink-0"
        :class="[
        post.postType === 'gallery' ? 'md:hidden' : '',
        showMetadata ? 'opacity-0 pointer-events-none md:opacity-100 md:pointer-events-auto' : 'opacity-100',
        post.postType === 'article' ? 'pb-4' : 'pb-12'
      ]"
    >
      <div
          class="flex items-center gap-4 text-[10px] font-bold tracking-widest uppercase text-muted-foreground/80 mb-6">
        <span>{{ formatDate(post.createdTime) }}</span>
        <span v-if="post.category" class="w-1 h-1 bg-border rounded-full"></span>
        <span v-if="post.category">{{ post.category.name }}</span>
      </div>
      <h1 v-if="post.postType !== 'article'"
          class="text-4xl md:text-6xl text-foreground font-normal leading-tight max-w-3xl">
        {{ post.title }}
      </h1>
      <p v-if="post.summary && post.postType === 'gallery'"
         class="text-muted-foreground text-xl leading-relaxed max-w-2xl font-light mt-8 max-h-[4.5em] overflow-y-auto scrollbar-hide">
        {{ post.summary }}
      </p>
    </header>

    <!-- Content: Internal Scroll Area -->
    <div
        class="flex-1 min-h-0 scrollbar-hide pb-32"
        :class="{
        'overflow-y-auto': post.postType === 'article' || !['gallery', 'video'].includes(post.postType),
        'overflow-hidden': ['gallery', 'video'].includes(post.postType)
      }"
    >
      <div v-if="post.postType === 'article'" class="max-w-[1024px] mx-auto px-8">
        <MarkdownRenderer :content="post.content"/>
      </div>
      <div v-else-if="post.postType === 'gallery'" class="h-full w-full relative">
        <!-- Desktop Layout: Side by Side -->
        <div class="hidden md:flex h-full max-w-[1400px] mx-auto">
          <!-- Photo Area (Left) - 60% -->
          <div class="w-[60%] h-full flex-shrink-0">
            <GalleryViewer
                :media-assets="post.mediaAssets"
                class="h-full"
                @active-change="handleActiveChange"
            />
          </div>

          <!-- Info Panel (Right) - 40% - Always Visible -->
          <div
              class="w-[40%] flex-shrink-0 overflow-y-auto scrollbar-hide bg-background border-l border-border/30 px-8 py-6 flex items-center">
            <div class="max-w-md mx-auto w-full">
              <!-- Toggle between Post Info and EXIF -->
              <Transition name="fade" mode="out-in">
                <!-- Post Info (Default) -->
                <div v-if="!showMetadata" key="post-info" class="space-y-8">
                  <!-- Post Header -->
                  <div class="space-y-4">
                    <div
                        class="flex items-center gap-3 text-[10px] font-bold tracking-widest uppercase text-muted-foreground/60">
                      <span>{{ formatDate(post.createdTime) }}</span>
                    </div>
                    <h1 class="text-3xl font-bold tracking-tight leading-tight">{{ post.title }}</h1>
                    <p v-if="post.summary" class="text-sm text-muted-foreground/80 leading-relaxed">
                      {{ post.summary }}
                    </p>
                  </div>
                </div>

                <!-- EXIF Info (When toggled) -->
                <div v-else key="exif-info" class="space-y-4 no-scrollbar">
                  <!-- Header -->
                  <div class="space-y-1 pb-3 border-b border-border/10">
                    <h3 class="text-xl font-bold tracking-tight">{{ $t('exif.title') }}</h3>
                    <p class="text-[11px] font-bold uppercase tracking-[0.2em] text-muted-foreground/30">
                      {{ currentAsset?.width }} × {{ currentAsset?.height }}
                    </p>
                  </div>

                  <!-- Exposure Section -->
                  <div class="space-y-4">
                    <h4 class="text-[10px] font-bold uppercase tracking-[0.3em] text-primary/60">
                      {{ $t('exif.sections.exposure') }}</h4>
                    <div class="grid grid-cols-2 gap-x-6 gap-y-4">
                      <div v-if="currentAsset?.metadata?.iso" class="space-y-1">
                        <span class="text-[12px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                            $t('exif.labels.iso')
                          }}</span>
                        <p class="text-[13px] font-mono font-bold">{{ currentAsset.metadata.iso }}</p>
                      </div>
                      <div v-if="currentAsset?.metadata?.fNumber" class="space-y-1">
                        <span class="text-[12px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                            $t('exif.labels.aperture')
                          }}</span>
                        <p class="text-[13px] font-mono font-bold">ƒ/{{ currentAsset.metadata.fNumber }}</p>
                      </div>
                      <div v-if="currentAsset?.metadata?.exposureTime" class="space-y-1">
                        <span class="text-[12px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                            $t('exif.labels.shutter')
                          }}</span>
                        <p class="text-[13px] font-mono font-bold">{{ currentAsset.metadata.exposureTime }}s</p>
                      </div>
                      <div v-if="currentAsset?.metadata?.exposureBias" class="space-y-1">
                        <span class="text-[12px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                            $t('exif.labels.expBias')
                          }}</span>
                        <p class="text-[13px] font-mono font-bold">{{ currentAsset.metadata.exposureBias }} EV</p>
                      </div>
                    </div>
                  </div>

                  <!-- Equipment Section -->
                  <div class="space-y-2.5">
                    <h4 class="text-[10px] font-bold uppercase tracking-[0.3em] text-primary/60">
                      {{ $t('exif.sections.equipment') }}</h4>
                    <div class="space-y-2.5">
                      <div v-if="currentAsset?.deviceMake || currentAsset?.deviceModel" class="space-y-1">
                        <span class="text-[12px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                            $t('exif.labels.camera')
                          }}</span>
                        <p class="text-[13px] font-bold leading-tight tracking-tight">{{ currentAsset.deviceMake }}
                          {{ currentAsset.deviceModel }}</p>
                      </div>
                      <div v-if="currentAsset?.metadata?.lensModel" class="space-y-1">
                        <span class="text-[12px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                            $t('exif.labels.lens')
                          }}</span>
                        <p class="text-[13px] font-bold leading-tight tracking-tight">{{
                            currentAsset.metadata.lensModel
                          }}</p>
                      </div>
                      <div v-if="currentAsset?.metadata?.focalLength" class="space-y-0.5">
                        <span class="text-[12px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                            $t('exif.labels.focalLength')
                          }}</span>
                        <p class="text-[13px] font-mono font-bold">{{ currentAsset.metadata.focalLength }}</p>
                      </div>
                    </div>
                  </div>

                  <!-- Processing Section -->
                  <div
                      v-if="currentAsset?.metadata && (getMeta('filmSimulation') || getMeta('filmMode') || getMeta('whiteBalance') || getMeta('dynamicRange') || getMeta('highlightTone') || getMeta('shadowTone') || getMeta('saturation') || getMeta('sharpness') || getMeta('noiseReduction') || getMeta('colorChromeEffect') || getMeta('colorChromeFXBlue') || getMeta('grainEffectRoughness'))"
                      class="space-y-4">
                    <h4 class="text-[10px] font-bold uppercase tracking-[0.3em] text-primary/60">
                      {{ $t('exif.sections.processing') }}</h4>
                    <div class="space-y-4">
                      <!-- Main Simulation -->
                      <div v-if="getMeta('filmSimulation') || getMeta('filmMode')" class="space-y-1">
                        <span class="text-[12px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                            $t('exif.labels.filmSimulation')
                          }}</span>
                        <p class="text-[13px] font-bold text-primary tracking-tight">
                          {{ tExifValue('filmSimulations', getMeta('filmSimulation') || getMeta('filmMode')) }}</p>
                      </div>

                      <!-- WB Details -->
                      <div class="grid grid-cols-2 gap-x-6">
                        <div v-if="getMeta('whiteBalance')" class="space-y-1">
                          <span class="text-[12px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                              $t('exif.labels.whiteBalance')
                            }}</span>
                          <p class="text-[13px] font-bold">{{ tExifValue('whiteBalance', getMeta('whiteBalance')) }}</p>
                        </div>
                        <div v-if="getMeta('whiteBalanceFineTune')" class="space-y-1">
                          <span class="text-[12px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                              $t('exif.labels.whiteBalanceFineTune')
                            }}</span>
                          <p class="text-[13px] font-mono font-bold">{{ getMeta('whiteBalanceFineTune') }}</p>
                        </div>
                      </div>

                      <!-- Recipe Grid - Balanced -->
                      <div class="pt-3 border-t border-border/5 grid grid-cols-3 gap-y-3 gap-x-2">
                        <div v-if="getMeta('dynamicRange')" class="space-y-0.5">
                          <span
                              class="text-[12px] font-bold uppercase tracking-widest text-muted-foreground/75 line-clamp-1 truncate">{{
                              $t('exif.labels.dynamicRange')
                            }}</span>
                          <p class="text-[13px] font-bold">{{ tExifValue('dynamicRange', getMeta('dynamicRange')) }}</p>
                        </div>
                        <div v-if="getMeta('highlightTone')" class="space-y-0.5">
                          <span class="text-[12px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                              $t('exif.labels.highlightTone')
                            }}</span>
                          <p class="text-[13px] font-bold">{{ tExifValue('generic', getMeta('highlightTone')) }}</p>
                        </div>
                        <div v-if="getMeta('shadowTone')" class="space-y-0.5">
                          <span class="text-[12px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                              $t('exif.labels.shadowTone')
                            }}</span>
                          <p class="text-[13px] font-bold">{{ tExifValue('generic', getMeta('shadowTone')) }}</p>
                        </div>
                        <div v-if="getMeta('saturation')" class="space-y-0.5">
                          <span class="text-[12px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                              $t('exif.labels.saturation')
                            }}</span>
                          <p class="text-[13px] font-bold">{{ tExifValue('generic', getMeta('saturation')) }}</p>
                        </div>
                        <div v-if="getMeta('sharpness')" class="space-y-0.5">
                          <span class="text-[12px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                              $t('exif.labels.sharpness')
                            }}</span>
                          <p class="text-[13px] font-bold">{{ tExifValue('generic', getMeta('sharpness')) }}</p>
                        </div>
                        <div v-if="getMeta('noiseReduction')" class="space-y-0.5">
                          <span class="text-[12px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                              $t('exif.labels.noiseReduction')
                            }}</span>
                          <p class="text-[13px] font-bold">{{ tExifValue('generic', getMeta('noiseReduction')) }}</p>
                        </div>
                      </div>

                      <!-- Effects - Cleaner -->
                      <div
                          v-if="getMeta('colorChromeEffect') || getMeta('colorChromeFXBlue') || getMeta('grainEffectRoughness')"
                          class="pt-2 flex flex-wrap gap-x-4 gap-y-2">
                        <div v-if="getMeta('colorChromeEffect')" class="flex items-center gap-2 grayscale opacity-70">
                          <span class="text-[12px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                              $t('exif.labels.colorChromeEffect')
                            }}</span>
                          <p class="text-[12px] font-bold">{{ tExifValue('generic', getMeta('colorChromeEffect')) }}</p>
                        </div>
                        <div v-if="getMeta('colorChromeFXBlue')" class="flex items-center gap-2 grayscale opacity-70">
                          <span class="text-[12px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                              $t('exif.labels.colorChromeFXBlue')
                            }}</span>
                          <p class="text-[12px] font-bold">{{ tExifValue('generic', getMeta('colorChromeFXBlue')) }}</p>
                        </div>
                        <div v-if="getMeta('grainEffectRoughness')"
                             class="flex items-center gap-2 grayscale opacity-70">
                          <span class="text-[12px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                              $t('exif.labels.grainEffectRoughness')
                            }}</span>
                          <p class="text-[12px] font-bold">{{
                              tExifValue('generic', getMeta('grainEffectRoughness'))
                            }}</p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </Transition>
            </div>
          </div>
        </div>

        <!-- Mobile Layout: Full Screen Photo + Bottom Panel -->
        <div class="md:hidden h-full">
          <GalleryViewer
              :media-assets="post.mediaAssets"
              class="h-full"
              @active-change="handleActiveChange"
          />

          <!-- Mobile Bottom Panel -->
          <Transition name="slide-up">
            <div v-if="showMetadata && currentAsset"
                 @touchmove.prevent
                 class="fixed inset-0 z-[60] bg-background/98 backdrop-blur-3xl overflow-hidden touch-none">
              <div class="h-full flex flex-col justify-center p-6 max-w-2xl mx-auto">
                <!-- Header -->
                <div class="flex items-start justify-between mb-4 shrink-0">
                  <div class="space-y-0.5">
                    <h3 class="text-lg font-bold">{{ $t('exif.mobileTitle') }}</h3>
                    <p class="text-[11px] font-bold uppercase tracking-[0.2em] text-muted-foreground/50">
                      {{ currentAsset.width }} × {{ currentAsset.height }}
                    </p>
                  </div>
                  <button @click="showMetadata = false"
                          class="p-1.5 hover:bg-muted rounded-full transition-colors opacity-50">
                    <X :size="18"/>
                  </button>
                </div>

                <!-- Content -->
                <div class="space-y-5">
                  <!-- Exposure -->
                  <div class="space-y-3">
                    <h4 class="text-[11px] font-bold uppercase tracking-[0.3em] text-primary/60">
                      {{ $t('exif.sections.exposure') }}</h4>
                    <div class="grid grid-cols-2 gap-4">
                      <div v-if="currentAsset.metadata?.iso" class="space-y-0.5">
                        <span class="text-[11px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                            $t('exif.labels.iso')
                          }}</span>
                        <p class="text-sm font-mono font-bold">{{ currentAsset.metadata.iso }}</p>
                      </div>
                      <div v-if="currentAsset.metadata?.fNumber" class="space-y-0.5">
                        <span class="text-[11px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                            $t('exif.labels.aperture')
                          }}</span>
                        <p class="text-sm font-mono font-bold">ƒ/{{ currentAsset.metadata.fNumber }}</p>
                      </div>
                      <div v-if="currentAsset.metadata?.exposureTime" class="space-y-0.5">
                        <span class="text-[11px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                            $t('exif.labels.shutter')
                          }}</span>
                        <p class="text-sm font-mono font-bold">{{ currentAsset.metadata.exposureTime }}s</p>
                      </div>
                      <div v-if="currentAsset.metadata?.exposureBias" class="space-y-0.5">
                        <span class="text-[11px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                            $t('exif.labels.expBias')
                          }}</span>
                        <p class="text-sm font-mono font-bold">{{ currentAsset.metadata.exposureBias }} EV</p>
                      </div>
                    </div>
                  </div>

                  <!-- Equipment -->
                  <div class="space-y-3">
                    <h4 class="text-[11px] font-bold uppercase tracking-[0.3em] text-primary/60">
                      {{ $t('exif.sections.equipment') }}</h4>
                    <div class="grid grid-cols-1 gap-3">
                      <div v-if="currentAsset.deviceMake || currentAsset.deviceModel" class="space-y-0.5">
                        <span class="text-[11px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                            $t('exif.labels.camera')
                          }}</span>
                        <p class="text-sm font-bold leading-tight">{{ currentAsset.deviceMake }}
                          {{ currentAsset.deviceModel }}</p>
                      </div>
                      <div v-if="currentAsset.metadata?.lensModel" class="space-y-0.5">
                        <span class="text-[11px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                            $t('exif.labels.lens')
                          }}</span>
                        <p class="text-sm font-bold leading-tight">{{ currentAsset.metadata.lensModel }}</p>
                      </div>
                      <div v-if="currentAsset.metadata?.focalLength" class="space-y-0.5">
                        <span class="text-[11px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                            $t('exif.labels.focalLength')
                          }}</span>
                        <p class="text-sm font-mono font-bold">{{ currentAsset.metadata.focalLength }}</p>
                      </div>
                    </div>
                  </div>

                  <!-- Processing -->
                  <div
                      v-if="currentAsset?.metadata && (getMeta('filmSimulation') || getMeta('filmMode') || getMeta('whiteBalance') || getMeta('dynamicRange') || getMeta('highlightTone') || getMeta('shadowTone') || getMeta('saturation') || getMeta('sharpness') || getMeta('noiseReduction') || getMeta('colorChromeEffect') || getMeta('colorChromeFXBlue') || getMeta('grainEffectRoughness'))"
                      class="space-y-5">
                    <h4 class="text-[11px] font-bold uppercase tracking-[0.3em] text-primary/60">
                      {{ $t('exif.sections.processing') }}</h4>
                    <div class="space-y-5">
                      <div v-if="getMeta('filmSimulation') || getMeta('filmMode')" class="space-y-0.5">
                        <span class="text-[11px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                            $t('exif.labels.filmSimulation')
                          }}</span>
                        <p class="text-sm font-bold text-primary tracking-tight">
                          {{ tExifValue('filmSimulations', getMeta('filmSimulation') || getMeta('filmMode')) }}</p>
                      </div>

                      <div class="grid grid-cols-2 gap-4">
                        <div v-if="getMeta('whiteBalance')" class="space-y-0.5">
                          <span class="text-[11px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                              $t('exif.labels.whiteBalance')
                            }}</span>
                          <p class="text-sm font-bold">{{ tExifValue('whiteBalance', getMeta('whiteBalance')) }}</p>
                        </div>
                        <div v-if="getMeta('whiteBalanceFineTune')" class="space-y-0.5">
                          <span class="text-[11px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                              $t('exif.labels.whiteBalanceFineTune')
                            }}</span>
                          <p class="text-sm font-mono font-bold">{{ getMeta('whiteBalanceFineTune') }}</p>
                        </div>
                      </div>

                      <!-- Mobile Recipe Grid - Readable -->
                      <div class="grid grid-cols-3 gap-y-4 gap-x-4 pt-4 border-t border-border/10">
                        <div v-if="getMeta('highlightTone')" class="space-y-0.5">
                          <span class="text-[11px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                              $t('exif.labels.highlightTone')
                            }}</span>
                          <p class="text-[11px] font-bold">{{ tExifValue('generic', getMeta('highlightTone')) }}</p>
                        </div>
                        <div v-if="getMeta('shadowTone')" class="space-y-0.5">
                          <span class="text-[11px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                              $t('exif.labels.shadowTone')
                            }}</span>
                          <p class="text-[11px] font-bold">{{ tExifValue('generic', getMeta('shadowTone')) }}</p>
                        </div>
                        <div v-if="getMeta('saturation')" class="space-y-0.5">
                          <span class="text-[11px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                              $t('exif.labels.saturation')
                            }}</span>
                          <p class="text-[11px] font-bold">{{ tExifValue('generic', getMeta('saturation')) }}</p>
                        </div>
                        <div v-if="getMeta('sharpness')" class="space-y-0.5">
                          <span class="text-[11px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                              $t('exif.labels.sharpness')
                            }}</span>
                          <p class="text-[11px] font-bold">{{ tExifValue('generic', getMeta('sharpness')) }}</p>
                        </div>
                        <div v-if="getMeta('noiseReduction')" class="space-y-0.5">
                          <span class="text-[11px] font-bold uppercase tracking-widest text-muted-foreground/75">{{
                              $t('exif.labels.noiseReduction')
                            }}</span>
                          <p class="text-[11px] font-bold">{{ tExifValue('generic', getMeta('noiseReduction')) }}</p>
                        </div>
                        <div v-if="getMeta('dynamicRange')" class="space-y-0.5">
                          <span
                              class="text-[11px] font-bold uppercase tracking-widest text-muted-foreground/75 line-clamp-1 truncate">{{
                              $t('exif.labels.dynamicRange')
                            }}</span>
                          <p class="text-[11px] font-bold">{{ tExifValue('dynamicRange', getMeta('dynamicRange')) }}</p>
                        </div>
                      </div>

                      <!-- Mobile Effects - Cleaner -->
                      <div
                          v-if="getMeta('colorChromeEffect') || getMeta('colorChromeFXBlue') || getMeta('grainEffectRoughness')"
                          class="pt-6 border-t border-border/10 flex flex-wrap gap-x-6 gap-y-4">
                        <div v-if="getMeta('colorChromeEffect')" class="space-y-1">
                          <span
                              class="text-[13px] font-bold uppercase tracking-widest text-muted-foreground/75 line-clamp-1 truncate">{{
                              $t('exif.labels.colorChromeEffect')
                            }}</span>
                          <p class="text-[13px] font-bold">{{ tExifValue('generic', getMeta('colorChromeEffect')) }}</p>
                        </div>
                        <div v-if="getMeta('colorChromeFXBlue')" class="space-y-1">
                          <span
                              class="text-[13px] font-bold uppercase tracking-widest text-muted-foreground/75 line-clamp-1 truncate">{{
                              $t('exif.labels.colorChromeFXBlue')
                            }}</span>
                          <p class="text-[13px] font-bold">{{ tExifValue('generic', getMeta('colorChromeFXBlue')) }}</p>
                        </div>
                        <div v-if="getMeta('grainEffectRoughness')" class="space-y-1">
                          <span
                              class="text-[13px] font-bold uppercase tracking-widest text-muted-foreground/75 line-clamp-1 truncate">{{
                              $t('exif.labels.grainEffectRoughness')
                            }}</span>
                          <p class="text-[13px] font-bold">{{
                              tExifValue('generic', getMeta('grainEffectRoughness'))
                            }}</p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </Transition>
        </div>
      </div>
      <div v-else-if="post.postType === 'video' && post.mediaAssets?.length"
           class="flex items-center justify-center h-full max-w-7xl mx-auto w-full md:w-[80%] lg:w-[60%] px-4">
        <MediaRenderer :id="post.mediaAssets[0].id" media-type="video" class="w-full h-full"/>
      </div>
      <div v-else class="max-w-4xl mx-auto px-4 text-lg text-foreground leading-relaxed whitespace-pre-wrap">
        {{ post.content }}
      </div>
    </div>

    <!-- Floating Immersive Action Bar -->
    <div class="fixed bottom-10 left-1/2 -translate-x-1/2 z-50">
      <div
          class="flex items-center gap-1.5 p-1.5 rounded-full bg-background/60 backdrop-blur-2xl border border-border/50 shadow-[0_20px_50px_-12px_rgba(0,0,0,0.2)] ring-1 ring-black/5">
        <!-- Navigation Group -->
        <button
            @click="goBack"
            class="flex items-center gap-2 pl-4 pr-5 py-2.5 rounded-full bg-foreground text-background font-bold text-sm hover:scale-105 transition-all active:scale-95 shadow-lg group whitespace-nowrap"
        >
          <ArrowLeft class="w-4 h-4 group-hover:-translate-x-0.5 transition-transform"/>
          <span>返回</span>
        </button>

        <!-- Media Info Toggle (Only for Gallery/Video) -->
        <template v-if="['gallery', 'video'].includes(post.postType)">
          <div class="w-px h-6 bg-border/40 mx-1"></div>
          <button
              @click="showMetadata = !showMetadata"
              class="p-2.5 rounded-full transition-all group relative"
              :class="showMetadata ? 'bg-primary text-primary-foreground' : 'hover:bg-muted text-muted-foreground hover:text-foreground'"
          >
            <Info :size="20" class="group-hover:rotate-12 transition-transform"/>
            <span v-if="showMetadata" class="absolute -top-1 -right-1 flex h-2 w-2">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2 w-2 bg-primary"></span>
            </span>
          </button>
        </template>
      </div>
    </div>

  </article>

  <div v-else-if="loading" class="py-24 space-y-8 animate-pulse max-w-4xl mx-auto px-4">
    <div class="h-4 bg-muted w-24 rounded"></div>
    <div class="h-10 bg-muted w-3/4 rounded-xl"></div>
    <div class="h-6 bg-muted w-1/2 rounded-lg"></div>
    <div class="h-96 bg-muted w-full pt-12 rounded-3xl"></div>
  </div>

  <div v-else class="min-h-screen flex flex-col items-center justify-center p-6 text-center space-y-8 pb-32">
    <div class="w-20 h-20 rounded-full bg-muted/40 flex items-center justify-center ring-1 ring-border/5">
      <HelpCircle class="w-10 h-10 text-muted-foreground/30 font-thin"/>
    </div>
    <div class="space-y-2">
      <h3 class="text-2xl font-bold tracking-tight">无法找到该文章</h3>
      <p class="text-sm text-muted-foreground/50 max-w-[280px] mx-auto leading-relaxed">
        它是如此之轻盈，以至于在光影中不留痕迹。
      </p>
    </div>
    <router-link
        to="/"
        class="flex items-center gap-2 px-8 py-3 rounded-full bg-foreground text-background text-xs font-bold uppercase tracking-widest hover:scale-105 active:scale-95 transition-all shadow-xl"
    >
      <ArrowLeft class="w-3.5 h-3.5"/>
      <span>{{ $t('nav.goHome') }}</span>
    </router-link>
  </div>

</template>

<script setup>
import {ref, onMounted} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {useI18n} from 'vue-i18n'
import dayjs from 'dayjs'

import {ArrowLeft, Sun, Moon, Monitor, Info, X, HelpCircle} from 'lucide-vue-next'
import {useUiStore} from '../stores/ui'
import api from '../api/client'
import MarkdownRenderer from '../components/MarkdownRenderer.vue'
import GalleryViewer from '../components/GalleryViewer.vue'
import MediaRenderer from '../components/MediaRenderer.vue'

const route = useRoute()
const router = useRouter()
const uiStore = useUiStore()
const {t, te, tm, locale} = useI18n()
const post = ref(null)
const loading = ref(true)
const showMetadata = ref(false)
const currentAsset = ref(null)

const getMeta = (key) => {
  if (!currentAsset.value?.metadata) return null
  const meta = currentAsset.value.metadata

  // 1. Exact match
  if (meta[key] !== undefined) return meta[key]

  // 2. Case-insensitive and snake_case match
  const searchKey = key.toLowerCase()
  const snakeKey = key.replace(/[A-Z]/g, letter => `_${letter.toLowerCase()}`)

  for (const k in meta) {
    const lk = k.toLowerCase()
    if (lk === searchKey || lk === snakeKey || lk === key.toLowerCase()) {
      return meta[k]
    }
  }
  return null
}

const tExifValue = (category, value) => {
  if (!value) return ''

  const valStr = String(value).trim()

  // Get all messages for the category
  const categoryMsgs = tm(`exif.values.${category}`)

  // Helper to find a match in a message object (Proxy or Plain)
  const findMatch = (msgs, v) => {
    if (!msgs || typeof msgs !== 'object') return null
    // Try exact match
    if (msgs[v]) return msgs[v]

    // Try fuzzy match (case-insensitive and clean characters)
    const searchV = v.toLowerCase().replace(/[\s\.\-\+_\/]/g, '')
    for (const key in msgs) {
      if (key.toLowerCase().replace(/[\s\.\-\+_\/]/g, '') === searchV) {
        return msgs[key]
      }
    }
    return null
  }

  // 1. Try category specific match
  let match = findMatch(categoryMsgs, valStr)
  if (match) return match

  // 2. Specialized handling for Dynamic Range (e.g., "DR200" or "200%")
  if (category === 'dynamicRange') {
    const numericPart = valStr.replace(/[^0-9]/g, '')
    if (numericPart) {
      match = findMatch(categoryMsgs, numericPart)
      if (match) return match
    }
  }

  // 3. Try generic values (Strong, Weak, etc.)
  match = findMatch(tm('exif.values.generic'), valStr)
  if (match) return match

  return valStr
}

const goBack = () => {
  router.back()
}

const handleActiveChange = ({asset}) => {
  currentAsset.value = asset
}


const fetchPost = async () => {
  loading.value = true
  try {
    const response = await api.get(`/posts/${route.params.id}`)
    post.value = response
    if (post.value.mediaAssets?.length > 0) {
      currentAsset.value = post.value.mediaAssets[0]
    }
  } catch (err) {

    console.error('Failed to fetch post:', err)
  } finally {
    loading.value = false
  }
}

const formatDate = (date) => {
  const currentLocale = locale.value === 'zh' ? 'zh-cn' : 'en'
  return dayjs(date).locale(currentLocale).format('LL')
}

const getCurrentIndex = () => {
  if (!currentAsset.value || !post.value?.mediaAssets) return 0
  return post.value.mediaAssets.findIndex(asset => asset.id === currentAsset.value.id)
}


onMounted(fetchPost)
</script>

<style scoped>
/* Unified Transitions */
.slide-up-enter-active, .slide-up-leave-active {
  transition: transform 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}

.slide-up-enter-from, .slide-up-leave-to {
  transform: translateY(100%);
}

.fade-enter-active, .fade-leave-active {
  transition: all 0.5s cubic-bezier(0.16, 1, 0.3, 1);
}

.fade-enter-from, .fade-leave-to {
  opacity: 0;
  transform: translateX(10px);
}

/* Scrollbar Hide Utility */
.no-scrollbar::-webkit-scrollbar {
  display: none;
}

.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.gallery-container {
  /* Using h-[calc(100vh-80px)] to ensure space for footer/tabs if any */
  height: 100vh;
}
</style>
